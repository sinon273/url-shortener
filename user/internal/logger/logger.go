package user_logger

import (
	"fmt"
	"log/slog"
	"os"
)

func Init(level string) (*slog.Logger, error) {
	var lvl slog.Level
	err := lvl.UnmarshalText([]byte(level))
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	traceHandler := &TraceHandler{next: jsonHandler}
	return slog.New(traceHandler), nil

}
