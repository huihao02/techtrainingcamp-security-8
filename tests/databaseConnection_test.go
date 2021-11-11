package tests

import (
	"example/app"
	"fmt"

	//"github.com/stretchr/testify/assert"
	"example/utils"
	"testing"
)

func TestInit(t *testing.T) {
	utils.Init()
	var db = utils.GetConnection()
	//db.Exec("insert into user values (?,?,?,?)", "qyz", "13661553037", "123321", "123123")
	var user []app.User
	name := "qyz"
	err := db.Select(&user, "select user_name,phone_number,password,mostly_used_device_id from user where user_name=?", name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)
}
