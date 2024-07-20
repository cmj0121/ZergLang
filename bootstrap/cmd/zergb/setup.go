package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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

	// set the verbose level
	switch a.Verbose {
	case 0:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	log.Debug().Msg("completed the setup of the logger")
}

// setup the output writer
func (a *Args) setupWriter(output string) (io.WriteCloser, error) {
	// get the output writer
	switch output {
	case "":
		log.Debug().Msg("using the standard output as the writer")
		return os.Stdout, nil
	default:
		return os.Create(a.Output)
	}
}
