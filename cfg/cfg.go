package cfg

import (
	"github.com/gokits/cfg"
	"github.com/gokits/cfg/source/file"
	"github.com/gokits/stdlogger"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	filesource *file.File
	meta       *cfg.ConfigMeta
	glogger    stdlogger.LeveledLogger
	gvalidator *validator.Validate
)

type Config struct {
	NsResourceCheckBL []ResourceNameGroup `json:"NsResourceCheckBL"` // namespace validate blacklist
}

type ResourceNameGroup struct {
	Name         string `json:"Name"`
	GroupVersion string `json:"GroupVersion"`
}

func (c *Config) PostSwap(oldptr interface{}) {
	glogger.Infof("config swapped: %+v", c)
}

func Init(configfile string, logger stdlogger.LeveledLogger) (err error) {
	glogger = logger
	if filesource, err = file.NewFileSource(configfile, file.WithLogger(logger)); err != nil {
		return
	}
	meta = cfg.NewConfigMeta(Config{}, filesource, cfg.WithLogger(logger))
	gvalidator = validator.New()
	go meta.Run()
	if err = meta.WaitSynced(); err != nil {
		filesource.Close()
		return
	}
	return
}

func Fini() {
	meta.Stop()
	filesource.Close()
}

func Get() *Config {
	return meta.Get().(*Config)
}
