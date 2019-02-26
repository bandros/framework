package router

import (
	"github.com/bandros/framework/sample/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(r *gin.Engine) {
	r.Static("/asset", "./asset")
	r.Static("/public", "./public")
	r.LoadHTMLGlob("./pages/**/*")

	r.NoRoute(error404)
	r.NoMethod(error404)

	r.GET("/", controller.Dashboard)

}

func error404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "error/404", gin.H{
		"title": "Error 404",
	})
}
