package app

import (
	"example/utils"
	"github.com/jmoiron/sqlx"
)

func VerifyPhoneNum(phone_num string, db *sqlx.DB)(status bool){
	
	//判断该手机号是否已经被注册
	var user []User
	_ = db.Select(&user, "SELECT user_name,phone_number,password,mostly_used_device_id from user where phone_number=?", phone_num)
	return len(user) != 0
}

func InsertIntoUser( userName, phoneNumber, password, mostlyUsedDeviceId string, db *sqlx.DB, registerTimes map[string]int) int {
	//检查该手机号的总注册次数是否超标
	if _, exist := registerTimes[phoneNumber]; exist {
		if times, _ := registerTimes[phoneNumber]; times >= utils.LIMIT_REGISTER_TIMES {
			return 2
		}
	}

	//预处理SQL语句，防止SQL注入
	stmt, _ := db.Prepare("INSERT INTO user(user_name, phone_number, password, mostly_used_device_id) VALUES (?, ?, ?, ?)")
	result, err := stmt.Exec(userName, phoneNumber, password, mostlyUsedDeviceId)

	//判断插入是否成功
	if err != nil{
		return 1
	}
	if aff, _ := result.RowsAffected(); aff != 1 {
		return 1
	}

	//增加该手机号的总共注册次数
	if _, exist := registerTimes[phoneNumber]; !exist {
		registerTimes[phoneNumber] = 1
	} else {
		registerTimes[phoneNumber]++
	}

	return 0
}

