package zerg

import (
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/rs/zerolog/log"
)

type Compiler struct {
	module *ir.Module
}

func NewCompiler() *Compiler {
	return &Compiler{
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

// run the compiler from the given source code
func (c *Compiler) ToIR(w io.Writer) error {
	c.run()

	// return the LLVM IR
	if w == nil {
		log.Info().Msg("write the LLVM IR to the stdout")
		w = os.Stdout
	}

	_, err := c.module.WriteTo(w)
	return err
}

// compile the source code
func (c *Compiler) run() {
	c.prologue()
	defer c.epilogue()

	// define the main function
	main := c.module.NewFunc("main", types.I32)
	builder := main.NewBlock("")
	builder.NewRet(constant.NewInt(types.I32, 0))
}
