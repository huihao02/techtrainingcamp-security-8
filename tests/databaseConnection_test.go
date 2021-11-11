package tests

import (
	//"github.com/stretchr/testify/assert"
	"example/utils"
	"testing"
)

func TestInit(t *testing.T) {
	utils.Init()
	var db = utils.GetConnection()
	db.Exec("insert into user values (?,?,?,?)", "qyz", "13661553037", "123321", "123123")

}
