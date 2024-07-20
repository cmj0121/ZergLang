package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Args struct {
	// The verbose level of the command.
	Verbose int `short:"v" type:"counter" help:"Set the verbose level of the command."`
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
	return nil
}

// setup everything before running the command
func (a *Args) prologue() {
	a.setupLogger()
}

// cleanup everything after running the command
func (a *Args) epilogue() {
	log.Info().Msg("finished the command ...")
}

// setup logger by the known settings
func (a *Args) setupLogger() {
	// setup the logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Debug().Msg("completed the setup of the logger")
}

func main() {
	args := New()
	if err := args.ParseAndRun(); err != nil {
		log.Error().Err(err).Msg("failed to run the command")
	}
}
