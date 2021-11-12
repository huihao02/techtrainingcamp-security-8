package main

import (
	"example/app"
	"example/utils"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	//某IP地址的连续登录错误次数
	var failTimes map[string]int
	failTimes = make(map[string]int)
	//某IP地址因连续登录错误次数超过规定值而被禁止账号密码登录的起始时间
	var bannedTime map[string]time.Time
	bannedTime = make(map[string]time.Time)
	//某手机号上次注册时间
	var lastTime = make(map[string]time.Time)
	//某手机号总共注册次数
	var registerTimes = make(map[string]int)
	//数据库连接初始化
	utils.Init()
	//获取数据库连接
	db := utils.GetConnection()

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
	r.POST("/login", func(c *gin.Context) {
		user_name := c.PostForm("UserName")
		password := c.PostForm("Password")
		IP := c.PostForm("IP")
		DeviceID := c.PostForm("DeviceID")

		_, banned := bannedTime[IP]
		//如果目前仍处于禁止登录状态，检测是否已经过了禁止登录时间
		if banned {
			banStartTime := bannedTime[IP]
			//若禁止登录时间已过
			if time.Now().After(banStartTime.Add(utils.BANNED_TIME)) {
				delete(bannedTime, IP)
				delete(failTimes, IP)
			} else {
				//否则不登录
				return
			}
		}

		//验证用户信息，并返回登录状态
		loginStatus := app.VerifyUser(user_name, password, IP, DeviceID, db, failTimes, bannedTime)
		var code int
		var message string
		var sessionID = utils.GetNewSessionId()
		//0:正常登录成功  1:需要滑块验证才能登录  2:需要等待一段时间才能用账号密码登录  3:需要手机号验证才能登录  4:登陆异常
		var decisionType int
		_, banned = bannedTime[IP]
		if banned {
			code = 1
			message = "登录错误次数过多！您将在5分钟内不可用账号密码登录！"
			decisionType = 2
		} else if loginStatus == 0 {
			code = 0
			message = "登录成功！"
			decisionType = 0
		} else if loginStatus == 1 {
			code = 1
			message = "用户名不存在，请重新输入！\n" + "剩余尝试次数：" + strconv.Itoa(utils.MAX_LOGIN_FAIL_TIME-failTimes[IP])
			decisionType = 1
		} else if loginStatus == 2 {
			code = 1
			message = "密码错误，请重新输入！\n" + "剩余尝试次数：" + strconv.Itoa(utils.MAX_LOGIN_FAIL_TIME-failTimes[IP])
			decisionType = 1
		} else if loginStatus == 3 {
			code = 0
			message = "检测到您未在常用登陆设备上登录，需要进行手机号验证"
			decisionType = 3
		} else {
			//状况外的返回值
			fmt.Println("登录状态返回有误！")
			code = 1
			message = "登录异常！建议使用手机验证码进行登录！"
			decisionType = 3
		}
		fmt.Println(message)
		c.JSON(http.StatusOK, gin.H{
			"Code":    code,
			"Message": message,
			"Data": gin.H{
				"SessionID":    sessionID,
				"ExpireTime":   utils.EXPIRE_TIME,
				"DecisionType": decisionType,
			},
		})
	})
	//注册
	r.POST("/register", func(c *gin.Context) {
		userName := c.PostForm("UserName")
		phoneNumber := c.PostForm("PhoneNumber")
		password := c.PostForm("Password")
		//ip不在表内，故并不使用
		//ip := c.PostForm("IP")
		deviceId := c.PostForm("DeviceID")

		//判断该手机号是否已被注册
		if app.VerifyPhoneNum(phoneNumber, db) {
			c.JSON(http.StatusOK, gin.H{
				"Code" : 1,
				"Message" : "该手机号已被注册",
			})
			return
		}
		var code int
		var message string
		var sessionID = utils.GetNewSessionId()
		var decisionType int

		if _, exist := lastTime[phoneNumber]; exist {
			last := lastTime[phoneNumber]
			//若禁止注册时间已过
			if time.Now().After(last.Add(utils.PHONENUM_REGISTER_INTERVAL)) {
				delete(bannedTime, phoneNumber)
			}
		}
		//执行插入操作，result为结果的状态： 0:执行成功 ， 1：语句执行失败， 2：注册次数过多
		var result int
		_, exist := lastTime[phoneNumber];
		if !exist{
			result = app.InsertIntoUser(userName, phoneNumber, password, deviceId, db, registerTimes)
		}
		if _, exist := lastTime[phoneNumber]; exist {
			//距离上次注册小于指定时间
			code = 1
			message = "注册过于频繁， 请等待一段时间后再次注册"
			decisionType = 2
		} else if result == 0 && registerTimes[phoneNumber] >= 5 {
			//注册次数达到一定数目，进行滑块验证
			code = 0
			message = "注册成功"
			decisionType = 1
			lastTime[phoneNumber] = time.Now()
		} else if result == 0 {
			//正常注册成功
			code = 0
			message = "注册成功"
			decisionType = 1
			lastTime[phoneNumber] = time.Now()
		}else if result == 1 {
			//注册执行插入语句出错
			code = 1
			message = "注册失败"
			decisionType = 0
		} else if result == 2 {
			//注册总次数达到限制
			code = 1
			message = "注册失败，注册总次数达到限制"
			decisionType = 3
		}


		c.JSON(http.StatusOK, gin.H{
			"Code" : code,
			"Message" : message,
			"Data": gin.H{
				"SessionID":    sessionID,
				"ExpireTime":   utils.EXPIRE_TIME,
				"DecisionType": decisionType,
			},
		})
	})
	r.Run(":9999")

}
