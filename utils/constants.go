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
