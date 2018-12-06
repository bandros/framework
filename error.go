package framework

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func ErrorJson(err string, c *gin.Context) {

	if err != "" {
		ipaddress := c.ClientIP()
		if ipaddress == os.Getenv("ipAddress") || os.Getenv("env") == "dev" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  err,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "Something error, please try again",
			})
		}

	}
}

func ErrorHtml(err string, c *gin.Context) {
	if err != "" {
		ipaddress := c.ClientIP()
		if ipaddress == os.Getenv("ipAddress") || os.Getenv("env") == "dev" {
			c.HTML(http.StatusBadRequest, "error/400", gin.H{
				"msg": err,
			})
		} else {
			c.HTML(http.StatusBadRequest, "error/400", gin.H{
				"msg": "Something error, please try again",
			})
		}

	}
}

func Error403(err string, c *gin.Context) {
	if err != "" {
		c.HTML(http.StatusForbidden, "error/400", gin.H{
			"msg": err,
		})

	}
}

func ErrorJson403(err string, c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"code": http.StatusForbidden,
		"msg":  err,
	})
}
