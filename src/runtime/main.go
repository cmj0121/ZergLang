package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
		kong.Name("zerg"),
		kong.Description("Zerg programming language runtime"),
		kong.UsageOnError(),
	)

	if cli.Version {
		fmt.Printf("zerg version %s\n", Version)
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

	// TODO: Implement lexer, parser, and evaluator
	_ = source

	return nil
}
