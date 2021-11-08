package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// func Session(keyPairs string) gin.HandlerFunc {
// 	store := SessionConfig()
// 	return sessions.Sessions(keyPairs, store)
// }

// func SessionConfig() sessions.Store {
// 	sessionMaxAge := 3600
// 	sessionSecret := "test"
// 	store := cookie.NewStore([]byte(sessionSecret))
// 	store.Options(sessions.Options{
// 		MaxAge: sessionMaxAge,
// 		Path: "/",
// 	})
// 	return store
// }

const (
	DefaultLen = 6
	CollectNum = 100
	Expiration = 10 * time.Minute
)

var captchaStore = captcha.NewMemoryStore(CollectNum, Expiration)

func CaptchaConfig() {
	captcha.SetCustomStore(captchaStore)
}

func CaptchaImage(c *gin.Context, length int) {
	w, h := 107, 36 // 图片大小
	captchaId := captcha.NewLen(length)
	session := sessions.Default(c)
	session.Set("captcha", captchaId)
	_ = session.Save()
	_ = Serve(c.Writer, c.Request, captchaId, ".png", "zh", false, w, h)
}

func CaptchaPhone(c *gin.Context) {
	// 如果手机号不存在返回错误
	captchaId := captcha.New()
	digits := captcha.RandomDigits(DefaultLen)
	captchaStore.Set(captchaId, digits)
	session := sessions.Default(c)
	session.Set("captcha", captchaId)
	_ = session.Save()
	fmt.Println(digits)
	c.JSON(http.StatusOK, gin.H{
		"Code":    0,
		"Message": "success",
		"Data": gin.H{
			"VerifyCode":   bytes.Runes(digits),
			"ExpireTime":   Expiration,
			"DecisionType": 0, // 风控不成功则为1
		},
	})
}

func CaptchaVerify(c *gin.Context, code string) bool {
	session := sessions.Default(c)
	if captchaId := session.Get("captcha"); captchaId != nil {
		session.Delete("captcha")
		_ = session.Save()
		if captcha.VerifyString(captchaId.(string), code) {
			return true
		} else {
			return false
		}
	}
	return false
}

func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		_ = captcha.WriteImage(&content, id, width, height)
	default:
		return captcha.ErrNotFound
	}
	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}
