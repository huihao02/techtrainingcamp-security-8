package app

import (
	"example/utils"
	"github.com/jmoiron/sqlx"
	"time"
)

type User struct {
	UserName           string `db:"user_name"`
	PhoneNumber        string `db:"phone_number"`
	Password           string `db:"password"`
	MostlyUsedDeviceId string `db:"mostly_used_device_id"`
}

//第一个返回值判断登录是否成功，0代表登录成功，1代表用户名出错,2代表密码出错，3代表未在常用设备上登录,
func VerifyUser(userName string, password string, IP string, deviceID string, db *sqlx.DB, failTImes map[string]int, bannedTime map[string]time.Time) int {
	var user []User
	_ = db.Select(&user, "select user_name,phone_number,password,mostly_used_device_id from user where user_name=?", userName)
	//用户名出错
	if len(user) == 0 {
		remainingAvailableTime(IP, failTImes, bannedTime)
		return 1
	}
	//密码出错
	if user[0].Password != password {
		remainingAvailableTime(IP, failTImes, bannedTime)
		return 2
	}
	//未在常用设备上登录
	if user[0].MostlyUsedDeviceId != deviceID {
		return 3
	}
	//登录成功
	return 0
}

func remainingAvailableTime(IP string, failTimes map[string]int, bannedTime map[string]time.Time) {
	failTime, ok := failTimes[IP]
	if !ok {
		failTimes[IP] = 1
	} else {
		failTimes[IP] = failTime + 1
	}
	if failTimes[IP] == utils.MAX_LOGIN_FAIL_TIME {
		bannedTime[IP] = time.Now()
	}
}
