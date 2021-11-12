package utils

import (
	"time"
)

//最大连续登录错误次数
var MAX_LOGIN_FAIL_TIME = 5

//因登录失败次数过多而被禁止登录的时长
var BANNED_TIME time.Duration = 5 * time.Minute

//session过期时间
var EXPIRE_TIME time.Duration = 3 * time.Hour


//一个手机号的最大注册次数
var LIMIT_REGISTER_TIMES = 10

//一个手机号注册后多久才能重新注册
var PHONENUM_REGISTER_INTERVAL = 24 * time.Hour
