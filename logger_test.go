package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"

	errorsx "errors"

	"github.com/pkg/errors"
)

func TestLog(t *testing.T) {
	Configure(slog.LevelDebug, true)

	slog.Debug("debug message")
	slog.Info("info message")
	slog.Warn("warning message")
	slog.Error("error message")

	slog.DebugContext(context.Background(), "debug message", slog.Any("level", "debug"))
	slog.InfoContext(context.Background(), "info mesage", slog.Any("level", "info"))
	slog.WarnContext(context.Background(), "warn message", slog.Any("level", "warn"))
	slog.ErrorContext(context.Background(), "error message", slog.Any("level", "error"))
}

func TestPrettyLogger(t *testing.T) {
	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level:       slog.LevelDebug,
			ReplaceAttr: ReplaceAttr,
		},
	}
	handler := NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)

	logger.Debug(
		"executing database query",
		slog.String("query", "SELECT * FROM users"),
	)

	logger.Info(
		"image upload successful",
		slog.String("image_id", "39ud88"),
	)

	logger.Warn(
		"storage is 90% full",
		slog.String("available_space", "900.1 MB"),
	)

	logger.Error(
		"An error occurred while processing the request",
		slog.Any("err", errors.New("errors error")),
	)
}

type TestLogValuer struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (t TestLogValuer) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", t.Name),
		slog.String("pass", "****"),
	)
}

var _ slog.LogValuer = (*TestLogValuer)(nil)

func TestLogValue(t *testing.T) {
	Configure(slog.LevelInfo, true)
	v := TestLogValuer{
		Name:     "tom",
		Password: "pass",
	}
	slog.Info("test log value", slog.Any("value", v.LogValue()))
}

func TestPkgErrors(t *testing.T) {
	err := errors.New("e1")
	cause := errors.Cause(err)
	if causer, ok := cause.(interface{ StackTrace() errors.StackTrace }); ok {
		stackTrace := causer.StackTrace()
		for _, frame := range stackTrace {
			source := fmt.Sprintf("%+s", frame)
			line := fmt.Sprintf("%d", frame)

			source = strings.Replace(source, "\n\t", " ", -1)
			ss := strings.Split(source, " ") // source string slice
			source = ss[0]
			function := ss[1]

			fmt.Println("source:", source)
			fmt.Println("line:", line)
			fmt.Println("func:", function)
		}
	}
}

func TestJsonHandlerError(t *testing.T) {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: ReplaceAttr,
	})
	logger := slog.New(h)
	ctx := context.Background()
	err := errors.New("something happened")
	logger.ErrorContext(ctx, "image uploaded", slog.Any("error", err))
	logger.Error("image uploaded", slog.Any("error", err))

	err2 := errorsx.New("err2")
	err3 := errorsx.Join(err, err2)
	logger.ErrorContext(ctx, "test error2", slog.Any("error", err3))
}
