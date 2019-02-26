package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard/index", gin.H{
		"title": "Bandros Framework",
	})
}
