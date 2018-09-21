package framework

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

type ConfigNew struct {
	Config []map[string]string `yaml:"config"`
}

type Init struct {
	Begin *gin.Engine
}

func(r *Init) Get(){
	ReloadConfig()
	if Config("env")!="dev" {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Begin = gin.Default()

	store := cookie.NewStore([]byte(Config("sessionKey")))
	r.Begin.Use(sessions.Sessions(Config("sessionName"), store))
	//config.Router(r)
	//r.Run(":"+os.Getenv("portHost"))

}

func(r *Init) Run()  {
	r.Begin.Run(":"+os.Getenv("portHost"))
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusPermanentRedirect)
}

func(r *Init) RunTls(domain ...string) error {
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	return autotls.Run(r.Begin, domain...)
}

func(r *Init) RunCert(cert,key string) {
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	http.ListenAndServeTLS(":443", cert,key,r.Begin)
}

func Config(key string)  string{
	return os.Getenv(key)
}

func ReloadConfig()  {
	var cfg ConfigNew
	file, _ := ioutil.ReadFile("config.yml")
	yaml.Unmarshal(file, &cfg)
	for i,v := range cfg.Config[0] {
		os.Setenv(i,v)
	}
}