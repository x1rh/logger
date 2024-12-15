package logger

import (
	"log/slog"
	"path/filepath"
	"github.com/mdobak/go-xerrors"

)

type StackFrame struct {
	Func   string `json:"func"`
	Source string `json:"source"`
	Line   int    `json:"line"`
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
	trace := xerrors.StackTrace(err)
	if len(trace) == 0 {
		return nil
	}

	frames := trace.Frames()
	s := make([]StackFrame, len(frames))
	for i, v := range frames {
		f := StackFrame{
			Source: filepath.Join(
				filepath.Base(filepath.Dir(v.File)),
				filepath.Base(v.File),
			),
			Func: filepath.Base(v.Function),
			Line: v.Line,
		}
		s[i] = f
	}
	return s
}

// fmtErr returns a slog.Value with keys `msg` and `trace`. If the error
// does not implement interface { StackTrace() errors.StackTrace }, the `trace`
// key is omitted.
func fmtErr(err error) slog.Value {
	// var groupValues []slog.Attr
	// groupValues = append(groupValues, slog.String("msg", err.Error()))
	// frames := marshalStack(err)
	// if frames != nil {
	// 	groupValues = append(groupValues,
	// 		slog.Any("trace", frames),
	// 	)
	// }
	// return slog.GroupValue(groupValues...)
	
	mp := make(map[string]any, 2)
	mp["error"] = err.Error()
	mp["trace"] = marshalStack(err)
	return (slog.AnyValue(mp))
}

