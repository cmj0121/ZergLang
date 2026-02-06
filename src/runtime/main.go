package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xrspace/zerglang/runtime/evaluator"
	"github.com/xrspace/zerglang/runtime/lexer"
	"github.com/xrspace/zerglang/runtime/parser"
)

var Version = "0.1.0"

type CLI struct {
	File    string `arg:"" optional:"" help:"Zerg source file to run (.zg)"`
	Verbose bool   `short:"v" help:"Enable verbose logging"`
	Version bool   `help:"Show version information"`
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("zerg-bootstrap"),
		kong.Description("Zerg bootstrap runtime interpreter"),
		kong.UsageOnError(),
	)

	if cli.Version {
		fmt.Printf("zerg-bootstrap version %s\n", Version)
		os.Exit(0)
	}

	if cli.File == "" {
		fmt.Fprintln(os.Stderr, "error: missing required argument: <file>")
		os.Exit(1)
	}

	setupLogging(cli.Verbose)

	if err := run(cli.File); err != nil {
		log.Error().Err(err).Msg("execution failed")
		ctx.Exit(1)
	}
}

func setupLogging(verbose bool) {
	level := zerolog.InfoLevel
	if verbose {
		level = zerolog.DebugLevel
	}

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(level).
		With().
		Timestamp().
		Logger()
}

func run(filename string) error {
	log.Debug().Str("file", filename).Msg("loading source file")

	source, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	log.Debug().Int("bytes", len(source)).Msg("source loaded")

	l := lexer.New(string(source))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, msg := range p.Errors() {
			log.Error().Msg(msg)
		}
		return fmt.Errorf("parsing failed with %d error(s)", len(p.Errors()))
	}

	log.Debug().Int("statements", len(program.Statements)).Msg("parsed program")

	// Set program args for sys.args()
	evaluator.SetProgramArgs(os.Args)

	env := evaluator.NewEnvironmentWithBuiltins()
	result := evaluator.Eval(program, env)

	if result != nil {
		if evaluator.IsError(result) {
			return fmt.Errorf("%s", result.Inspect())
		}
		log.Debug().Str("result", result.Inspect()).Msg("evaluation complete")
	}

	return nil
}
