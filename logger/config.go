package logger

import (
	"io"
	"os"

	"github.com/nagokos/connefut_backend/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger

func init() {
	logPath := config.Config.LogFile

	logfile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("file=logFile err=%s", err.Error())
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: "2006-01-02 15:04:05"}
	multiLogFile := io.MultiWriter(output, logfile)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = output.TimeFormat

	Log = zerolog.New(multiLogFile).With().Timestamp().Caller().Logger()
}
