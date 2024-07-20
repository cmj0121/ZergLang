package main

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zerglang/bootstrap/pkg/zerg"
)

type Args struct {
	// The verbose level of the command.
	Verbose int `short:"v" type:"counter" help:"Set the verbose level of the command."`

	// the output options
	Output string `short:"o" name:"output" type:"path" help:"The output file to save the result."`
}

// create a new instance of Args with the default settings
func New() *Args {
	return &Args{}
}

func (a *Args) ParseAndRun() error {
	ctx := kong.Parse(a)
	return a.Run(ctx)
}

func (a *Args) Run(ctx *kong.Context) error {
	a.prologue()
	defer a.epilogue()

	log.Info().Msg("starting the command ...")
	return a.run()
}

func (a *Args) run() error {
	compiler := zerg.NewCompiler()

	writer, err := a.setupWriter(a.Output)
	if err != nil {
		log.Error().Err(err).Msg("failed to setup the writer")
		return err
	}
	defer writer.Close()

	return compiler.ToIR(writer)
}

func main() {
	args := New()
	if err := args.ParseAndRun(); err != nil {
		log.Error().Err(err).Msg("failed to run the command")
	}
}
