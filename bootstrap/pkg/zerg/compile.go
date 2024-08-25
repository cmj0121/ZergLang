package zerg

import (
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zerglang/bootstrap/pkg/zerg/lexer"
)

type Compiler struct {
	*lexer.Lexer

	module *ir.Module
}

func NewCompiler(r io.Reader) *Compiler {
	return &Compiler{
		Lexer:  lexer.New(r),
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

	for range c.Iterate(ctx) {
		//
	}

	if c.Err() != nil {
		log.Warn().Err(c.Err()).Msg("failed to compile the source code")
		return c.Err()
	}

	// define the main function
	main := c.module.NewFunc("main", types.I32)
	builder := main.NewBlock("")
	builder.NewRet(constant.NewInt(types.I32, 0))

	return nil
}
