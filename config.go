package framework

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
		r.Begin = gin.New()
	} else {
		r.Begin = gin.Default()
	}
	store := cookie.NewStore([]byte(Config("sessionKey")))
	r.Begin.Use(sessions.Sessions(Config("sessionName"), store))
}

func (r *Init) Run() {
	r.Begin.Run(":" + os.Getenv("portHost"))
}

func Config(key string) string {
	return os.Getenv(key)
}

func ReloadConfig() {
	var cfg ConfigNew
	var fileName = "config.yml"
	var exist = FileExist(fileName)
	if !exist{
		fileName = "config.yaml"
		exist = FileExist(fileName)
		if !exist{
			panic("config.yaml or config.yml doesn't exist")
		}
	}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	yaml.Unmarshal(file, &cfg)
	if len(cfg.Config) >= 1 {
		for i, v := range cfg.Config[0] {
			os.Setenv(i, v)
		}
	}

}
