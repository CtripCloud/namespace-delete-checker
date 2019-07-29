package cfg

import (
	"os"

	"github.com/gin-gonic/gin/json"
)

var config = &Cfg{}

type ResourceNameGroup struct {
	Name         string `json:"Name"`
	GroupVersion string `json:"GroupVersion"`
}

type Cfg struct {
	NsResourceCheckBL []ResourceNameGroup `json:"NsResourceCheckBL"` // namespace validate blacklist
}

// Config returns current configuration, should be the only entrance to get cfg
func Config() *Cfg {
	return config
}

// MustInit panic if fail
func MustInit(path string) {
	var file *os.File
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	if err = json.NewDecoder(file).Decode(config); err != nil {
		panic(err.Error())
	}
}
