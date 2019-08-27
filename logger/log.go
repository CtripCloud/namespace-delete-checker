package logger

import (
	"io"
	"path/filepath"

	"github.com/YueHonghui/rfw"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	// LogDir is where to store access/runtime logs
	LogDir string
	// LogRemain is days to remain logs in LogDir
	LogRemain int
	// Writers to write logs to
	access, runtime io.Writer
)

const (
	accessLogFile  = "access"
	runtimeLogFile = "runtime"
)

// NewOutputOrDie returns a rfw writer to write logs to
func NewOutputOrDie(file string) *rfw.Rfw {
	rfw, err := rfw.NewWithOptions(filepath.Join(LogDir, file), rfw.WithCleanUp(LogRemain))
	if err != nil {
		panic(err)
	}
	return rfw
}

// SetupAccessLog makes sure gin access logs be written to w
func SetupAccessLog(w io.Writer) {
	gin.DefaultWriter = w
	gin.DefaultErrorWriter = w
}

// SetupRuntimeLog setup runtime log
func SetupRuntimeLog(w io.Writer) {
	logrus.SetOutput(w)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
}

// MustInit setup access/runtime logs or die
func MustInit() {
	access = NewOutputOrDie(accessLogFile)
	SetupAccessLog(access)

	runtime = NewOutputOrDie(runtimeLogFile)
	SetupRuntimeLog(runtime)
}

// Close closes access/runtime writers, should be defer called in main
func Close() {
	if access != nil {
		access.(*rfw.Rfw).Close()
	}
	if runtime != nil {
		runtime.(*rfw.Rfw).Close()
	}
}
