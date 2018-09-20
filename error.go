package framework

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)


func ErrorJson(err string,c *gin.Context,){

	if err != "" {
		ipaddress := c.ClientIP()
		if(ipaddress==os.Getenv("ipAddress") || os.Getenv("env")=="dev"){
			c.JSON(http.StatusInternalServerError,gin.H{
				"code" : http.StatusInternalServerError,
				"msg" : err,
			})
		}else{
			c.JSON(http.StatusInternalServerError,gin.H{
				"code" : http.StatusInternalServerError,
				"msg" : "Something error, please try again",
			})
		}

	}
}

func ErrorHtml(err string,c *gin.Context){
	if err != "" {
		ipaddress := c.ClientIP()
		if(ipaddress==os.Getenv("ipAddress") || os.Getenv("env")=="dev"){
			c.HTML(http.StatusInternalServerError,"error/400",gin.H{
				"msg" : err,
			})
		}else{
			c.HTML(http.StatusInternalServerError,"error/400",gin.H{
				"msg" : "Something error, please try again",
			})
		}

	}
}
