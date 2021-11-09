package main

import (
	"example/app"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 设置验证码的自定义 store
	app.CaptchaConfig()
	// 获取图片验证码
	r.GET("/captcha/image", func(c *gin.Context) {
		app.CaptchaImage(c, 4)
	})
	// 获取手机验证码
	r.POST("/captcha/phone", func(c *gin.Context) {
		app.CaptchaPhone(c)
	})
	// 验证验证码 Code 0 成功，1 失败
	r.GET("/captcha/verify/:value", func(c *gin.Context) {
		value := c.Param("value")
		if app.CaptchaVerify(c, value) {
			c.JSON(http.StatusOK, gin.H{
				"Code":    0,
				"Message": "success",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"Code":    1,
				"Message": "failed",
			})
		}
	})
	r.Run(":9999")

}
