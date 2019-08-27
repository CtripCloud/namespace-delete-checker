# Stdlogger
Logger interface and wrappers for common implementations

# What this package provides
## Minimal leveled logger interface
A minimal leveled logger interface is provided. We can use this interface in other packages so the package would not be tied together with any specific logger implementation. 
```go
type LeveledLogger interface {
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
}
```

## Wrappers of common logger implementations
This package already includes some wrappers of logger implementations. And many other logger implementations already implemented `LeveledLogger` interface(`zap.Sugar` for example). If you want, you can write your own wrappers very easy.

example for [logrus](https://github.com/sirupsen/logrus)
```go
import "github.com/gokits/stdlogger/logrus"

//...

var logglobal LeveledLogger = logrus.FromGlobal()
var logentry LeveledLogger = logrus.FromEntry(...)
```

example for [zap](https://github.com/uber-go/zap)

```go
import "go.uber.org/zap"

//...

logger, _ := zap.NewProduction()
defer logger.Sync() // flushes buffer, if any
var sugar LeveledLogger = logger.Sugar()
```

# Status of this package
We have already used this package in our product environment. But this package has not been released version v1.0.0, so compatibility is not guaranteed through changes.

# How to contribute
Welcome to submit PRs!