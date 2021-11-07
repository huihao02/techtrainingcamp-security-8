package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)



func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 获取图片验证码
	r.GET("/captcha/image", func(c *gin.Context) {
		Captcha(c, 4)
	})
	// Code 0 成功，1 失败
	r.GET("/captcha/verify/:value", func(c *gin.Context) {
		value := c.Param("value")
		if CaptchaVerify(c, value) {
			c.JSON(http.StatusOK, gin.H{
				"Code": 0,
				"Message": "success",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"Code": 1,
				"Message": "failed",
			})
		}
	})
	r.Run(":9999")
}