package logger

import (
	"log/slog"
	// "path/filepath"

	// "github.com/mdobak/go-xerrors"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type StackFrame struct {
	Func   string `json:"func"`
	Source string `json:"source"`
	Line   string    `json:"line"`
}

// ReplaceAttr using for slog.Options.Replace
func ReplaceAttr(_ []string, attr slog.Attr) slog.Attr {
	switch attr.Value.Kind() {
	case slog.KindAny:
		switch v := attr.Value.Any().(type) {
		case error:
			attr.Value = fmtErr(v)
		}
	}
	return attr
}

// marshalStack extracts stack frames from the error
func marshalStack(err error) []StackFrame {
	if causer, ok := errors.Cause(err).(interface{ StackTrace() errors.StackTrace }); ok {
		stackTrace := causer.StackTrace()
		frames := make([]StackFrame, len(stackTrace))
		for i, frame := range stackTrace {
			source := fmt.Sprintf("%+s", frame)
			line := fmt.Sprintf("%d", frame)

			// github.com/pkg/errors is archived and the implementation for getting function name or source file is inconvenient
			// TODO: replace github.com/pkg/errors with github.com/mdobak/go-xerrors
			source = strings.Replace(source, "\n\t", " ", -1)
			ss := strings.Split(source, " ") // source string slice
			// source = ss[0]
			// function := ss[1]
			
			frames[i].Source = ss[0]
			frames[i].Func = ss[1]
			frames[i].Line = line 
		}
		return frames
	}

	return nil 
}


func fmtErr(err error) slog.Value {
	mp := make(map[string]any, 2)
	mp["message"] = err.Error()
	frames := marshalStack(err)
	if frames != nil {
		mp["trace"] = frames
	}
	return (slog.AnyValue(mp))
}
