package framework

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"
)

type AuthConfig struct {
	Role map[string]interface{} `yaml:"auth"`
}

func Auth(c *gin.Context)  {
	var auth AuthConfig
	file, _ := ioutil.ReadFile("auth.yaml")
	yaml.Unmarshal(file, &auth)
	token,err := getToken(c)
	if err != nil {
		Error403(err.Error(),c)
		c.Abort()
		return
	}
	role := token["role"]
	if role != nil &&  auth.Role[role.(string)] != nil{
		allow := auth.Role[role.(string)].([]interface{})
		show := false
		url := c.Request.URL.Path
		for _,v := range allow {
			pattern := v.(string)
			pattern = strings.Replace(pattern,"*","(.*)",-1)
			if r,_:=regexp.MatchString("^"+pattern+"$", url);r  {
				show = true
				fmt.Println("cocok : ","^"+v.(string)+"$",":",url)
				break
			}
			fmt.Println("belum ada")
		}
		if !show {
			Error403("Not found access",c)
			c.Abort()
			return
		} else {
			c.Set("jwt", token)
		}



	}else{
		Error403("Role not found",c)
		c.Abort()
		return
	}
}

func getToken(c *gin.Context) (jwt.MapClaims,error) {
	var tokenString string
	session := sessions.Default(c)
	v := session.Get(Config("jwtName"))
	if v == nil {
		return nil,errors.New("Error token")
	} else {
		tokenString = v.(string)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(Config("jwtKey")), nil
	})
	if token == nil || err != nil {
		return nil,errors.New("Error token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil,errors.New("Error token")
	}
	return claims,nil
}
