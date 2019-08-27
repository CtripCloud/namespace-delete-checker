package nooplogger

type NoopLogger struct{}

func (nl *NoopLogger) Debugf(format string, v ...interface{}) {
}

func (nl *NoopLogger) Debug(v ...interface{}) {
}

func (nl *NoopLogger) Infof(format string, v ...interface{}) {
}

func (nl *NoopLogger) Info(v ...interface{}) {
}

func (nl *NoopLogger) Warnf(format string, v ...interface{}) {
}

func (nl *NoopLogger) Warn(v ...interface{}) {
}

func (nl *NoopLogger) Errorf(format string, v ...interface{}) {
}

func (nl *NoopLogger) Error(v ...interface{}) {
}

var gnoop *NoopLogger = new(NoopLogger)

func Default() *NoopLogger {
	return gnoop
}
