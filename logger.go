package logger

import (
	"log/slog"
	"os"

)

var defaultLogLevel *slog.LevelVar
var defaultLogger *slog.Logger
var defaultAddSource bool = true

func init() {
	setDefaultLogLevel(slog.LevelDebug)
	Configure(defaultLogLevel, defaultAddSource)
}

// Configure() we only configure the default slog logger
// if you want a specific slog logger, just new one using package `log/slog`, instead this package
func Configure(level slog.Leveler, AddSource bool) {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: AddSource,

		// NOTICE: performance 
		ReplaceAttr: ReplaceAttr,
	}

	var handler slog.Handler
	switch level.Level() {
	case slog.LevelDebug:
		handler = NewPrettyHandler(os.Stdout, PrettyHandlerOptions{SlogOpts: *opts})
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

func setDefaultLogLevel(level slog.Level) {
	if defaultLogLevel == nil {
		defaultLogLevel = &slog.LevelVar{}
	}
	defaultLogLevel.Set(level)
}

func SetLogLevel(level slog.Level) {
	setDefaultLogLevel(level)
	Configure(level, defaultAddSource)
}

