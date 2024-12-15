## logger
`logger` is a simple and sufficiently practical logger package, using `log/slog`, aiming to avoid dependency on third-party libraries.

usage 
```go
import (
    "log/slog"
    _ "github.com/x1rh/logger"
)

func main() {
    // initialize it by importing it
    slog.Info("hello")
    slog.Info("message 1", slog.String("key", "value"))
    slog.Error("error")

    // or Configure() it again 
    logLevel := slog.LevelDebug
    addSource := true
    logger.Configure(logLevel, addSource)
}
```

Initialize the default slog logger once in the function main() and use it anywhere.
