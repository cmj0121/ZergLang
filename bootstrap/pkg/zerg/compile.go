package zerg

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zerglang/bootstrap/pkg/zerg/parser"
	"github.com/cmj0121/zerglang/bootstrap/pkg/zerg/token"
)

// The context key that hold the LLVM context
type ContextKey string

var (
	BuilderKey    ContextKey = "builder"
	ReturnTypeKey ContextKey = "return_type"
	ReturnedKey   ContextKey = "returned"
	ValueKey	  ContextKey = "value"
)

type Compiler struct {
	*parser.Parser

	module *ir.Module
}

func NewCompiler(r io.Reader) *Compiler {
	return &Compiler{
		Parser: parser.New(r),
		module: ir.NewModule(),
	}
}

func (c *Compiler) prologue() {
	c.setupTargetTriple()
}

func (c *Compiler) epilogue() {
}

func (c *Compiler) setupTargetTriple() {
	triple := []string{}

	// define the target architecture
	switch runtime.GOARCH {
	case "amd64":
		triple = append(triple, "x86_64")
	case "386":
		triple = append(triple, "i386")
	case "arm64":
		triple = append(triple, "aarch64")
	case "arm":
		triple = append(triple, "arm")
	default:
		triple = append(triple, "unknown")
	}

	// define the target triple
	switch runtime.GOOS {
	case "darwin":
		triple = append(triple, "apple-darwin")
	case "linux":
		triple = append(triple, "pc-linux")
	case "windows":
		triple = append(triple, "pc-windows")
	default:
		triple = append(triple, "unknown")
	}

	c.module.TargetTriple = strings.Join(triple, "-")
}

// run the compiler from the given source code to the object file
func (c *Compiler) ToIR(ctx context.Context, output string) error {
	if err := c.run(ctx); err != nil {
		log.Error().Err(err).Msg("failed to compile the source code")
		return err
	}

	var w io.WriteCloser
	switch output {
	case "":
		w = os.Stdout
	default:
		file, err := os.Create(output)
		if err != nil {
			log.Error().Err(err).Str("output", output).Msg("failed to create the output file")
			return err
		}
		defer file.Close()

		w = file
	}

	_, err := c.module.WriteTo(w)
	return err
}

// run the compiler from the given source code to the object file
func (c *Compiler) ToObj(ctx context.Context, output string) error {
	if err := c.run(ctx); err != nil {
		log.Error().Err(err).Msg("failed to compile the source code")
		return err
	}

	return c.buildTo(output, "-Wno-override-module", "-c")
}

// run the compiler from the given source code to the binary file
func (c *Compiler) ToBin(ctx context.Context, output string) error {
	if err := c.run(ctx); err != nil {
		log.Error().Err(err).Msg("failed to compile the source code")
		return err
	}

	return c.buildTo(output, "-Wno-override-module")
}

func (c *Compiler) buildTo(output string, args ...string) error {
	ir, err := os.CreateTemp("", "zergb-********.ll")
	if err != nil {
		log.Error().Err(err).Msg("failed to create the temporary file")
		return err
	}
	defer ir.Close()
	defer os.Remove(ir.Name())

	if _, err := c.module.WriteTo(ir); err != nil {
		log.Error().Err(err).Msg("failed to write the LLVM IR to the temporary file")
		return err
	}

	args = append(args, ir.Name())
	args = append(args, "-o", output)
	cmd := exec.Command("clang", args...)

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Str("output", string(stdout)).Msg("failed to build the object file")
		return err
	}

	return nil
}

// compile the source code
func (c *Compiler) run(ctx context.Context) error {
	c.prologue()
	defer c.epilogue()

	if err := c.Parse(ctx); err != nil {
		log.Warn().Err(err).Msg("failed to parse the source code")
		return err
	}

	if c.Err() != nil {
		log.Warn().Err(c.Err()).Msg("failed to compile the source code")
		return c.Err()
	}

	return c.compileAST(ctx, c.Root())
}

// Compile the AST to the LLVM IR
func (c *Compiler) compileAST(ctx context.Context, node *parser.Node) error {
	switch node.Type() {
	case parser.Root:
		for _, child := range node.Children() {
			if err := c.compileAST(ctx, child); err != nil {
				log.Warn().Err(err).Msg("failed to compile the child node")
				return err
			}
		}
	case parser.Fn:
		return c.compileFunction(ctx, node)
	case parser.PrintStmt:
		return c.compilePrintStmt(ctx, node)
	case parser.ReturnStmt:
		expr := node.Children()[0]
		if err := c.compileAST(ctx, expr); err != nil {
			log.Warn().Err(err).Msg("failed to compile the expression")
			return err
		}
	default:
		log.Warn().Any("type", node.Type()).Msg("unknown node type")
		return fmt.Errorf("unknown node type: %v", node.Type())
	}

	return nil
}

// Compile function and related statements.
func (c *Compiler) compileFunction(ctx context.Context, node *parser.Node) error {
	name := node.Token().String()

	_, type_hint, stmts := node.Children()[0], node.Children()[1], node.Children()[2]

	typ := c.toLLVMType(type_hint)

	fn := c.module.NewFunc(name, typ)
	builder := fn.NewBlock("")

	// set the builder to the context
	ctx = context.WithValue(ctx, BuilderKey, builder)
	ctx = context.WithValue(ctx, ReturnTypeKey, typ)

	for _, child := range stmts.Children() {
		if err := c.compileAST(ctx, child); err != nil {
			log.Warn().Err(err).Msg("failed to compile the child node")
			return err
		}
	}

	// check the return statement called or not
	if _, ok := ctx.Value(ReturnedKey).(bool); !ok {
		log.Debug().Msg("function does not have the return statement, add the default return")
		if err := c.compileReturn(ctx, nil); err != nil {
			log.Warn().Err(err).Msg("failed to add the default return statement")
			return err
		}
	}

	return nil
}

// Compile the print statement
func (c *Compiler) compilePrintStmt(ctx context.Context, node *parser.Node) error {
	expr := node.Children()[0]
	value, err := c.compileExpression(ctx, expr)
	if err != nil {
		log.Warn().Err(err).Msg("failed to compile the expression")
		return err
	}

	typ := value.Type();
	switch typ := typ.(type) {
	case *types.PointerType:
	default:
		err := fmt.Errorf("unsupport print type: %T", typ)
		log.Warn().Err(err).Any("type", typ).Msg("failed to print the value")
		return err
	}

	return c.showString(ctx, value)
}

// Compile the current function tnat return the zero value of the return type.
func (c *Compiler) compileReturn(ctx context.Context, value any) error {
	builder := ctx.Value(BuilderKey).(*ir.Block)
	typ := ctx.Value(ReturnTypeKey).(types.Type)

	v, err := c.toLLVMValue(typ, value)
	if err != nil {
		log.Warn().Err(err).Msg("failed to convert the value to the LLVM value")
		return err
	}

	builder.NewRet(v)
	return nil
}

// Compile the expression
func (c *Compiler) compileExpression(ctx context.Context, node *parser.Node) (value.Value, error) {
	switch typ := node.Token().Type(); typ {
	case token.Int:
		raw := node.Token().String()
		v, err := strconv.Atoi(raw)
		if err != nil {
			log.Warn().Err(err).Str("raw", raw).Msg("failed to convert the raw token to the integer")
			return nil, err
		}

		return c.toLLVMValue(types.I32, v)
	case token.String:
		raw := node.Token().String()
		str := constant.NewCharArrayFromString(raw)

		// create the global variable
		global := c.module.NewGlobalDef("", str)
		return global, nil
	default:
		err := fmt.Errorf("unknown token type: %v", typ)
		log.Warn().Err(err).Msg("failed to compile the expression")
		return nil, err
	}
}

// Show the string to the STDOUT
func (c *Compiler) showString(ctx context.Context, v value.Value) error {
	log.Info().Str("value", v.Ident()).Msg("show the string")
	return nil
}

// Get the LLVM IR type from the AST type hint
func (c *Compiler) toLLVMType(node *parser.Node) types.Type {
	switch node.Type() {
	case parser.Type:
		if node.Token() == nil {
			log.Info().Msg("no type hint, use the default type")
			return types.Void
		}

		switch node.Token().String() {
		case "u32":
			return types.I32
		case "void":
			return types.Void
		default:
			log.Warn().Str("type", node.Token().String()).Msg("unknown type hint")
			return types.Void
		}
	default:
		log.Warn().Any("type", node.Type()).Msg("unknown node type")
		return types.Void
	}
}

// Get the LLVM value from the passed value
func (c *Compiler) toLLVMValue(typ types.Type, v any) (value.Value, error) {
	switch typ := typ.(type) {
	case *types.VoidType:
		// always return the void type
		return nil, nil
	case *types.IntType:
		switch v := v.(type) {
		case int:
			return constant.NewInt(typ, int64(v)), nil
		default:
			log.Warn().Any("value", v).Msg("unknown value type")
			return nil, fmt.Errorf("unknown value type: %v", v)
		}
	default:
		log.Warn().Any("type", typ).Msg("unsupport LLVM type")
		return nil, fmt.Errorf("unsupport LLVM type: %v", typ)
	}
}
