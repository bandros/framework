package framework

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type ConfigNew struct {
	Config []map[string]string `yaml:"config"`
}

type Init struct {
	Begin *gin.Engine
}

func (r *Init) Get() {
	ReloadConfig()
	if Config("env") != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Begin = gin.Default()
}

func (r *Init) Run() {
	r.Begin.Run(":" + os.Getenv("portHost"))
}

func Config(key string) string {
	return os.Getenv(key)
}

func ReloadConfig() {
	var cfg ConfigNew
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		file,err = ioutil.ReadFile("config.yaml")
		if err != nil {
			panic(err.Error())
		}
	}
	yaml.Unmarshal(file, &cfg)
	if len(cfg.Config) >= 1 {
		for i, v := range cfg.Config[0] {
			os.Setenv(i, v)
		}
	}

}
