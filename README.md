## logger
`logger` is a simple and sufficiently practical logger package, using `log/slog`, aiming to avoid dependency on third-party libraries.

usage 
```go
import (
    "log/slog"
    "github.com/x1rh/logger"
)

func main() {
    logLevel := slog.LevelDebug
    addSource := true
    logger.Configure(logLevel, addSource)

    slog.Info("hello")
    slog.Error("error")
}
```

Initialize the default slog logger once in the function main() and use it anywhere.