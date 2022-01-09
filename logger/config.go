package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/nagokos/connefut_backend/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
)

type Logger struct{}

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

func (l *Logger) Log(ctx context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	if string(level) == "error" {
		w := bytes.NewBufferString(fmt.Sprintf("%s:%v", level, msg))

		for k, v := range data {
			fmt.Fprintf(w, "\t%s:%v", k, v)
		}

		fmt.Fprintf(w, "\x1b[49m")
		Log.Error().Msg(w.String())
	} else {
		w := bytes.NewBufferString(fmt.Sprintf("%s:%v", level, msg))

		for k, v := range data {
			fmt.Fprintf(w, "\t%s:%v", k, v)
		}

		fmt.Fprintf(w, "\x1b[49m")
		Log.Debug().Msg(w.String())
	}
}
