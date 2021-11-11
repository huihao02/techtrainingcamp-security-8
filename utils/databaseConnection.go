package utils

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

const (
	USER_NAME = "root"
	PASSWORD  = "18901785998qQ"
	HOST      = "localhost"
	PORT      = "3306"
	DATABASE  = "tech_training_camp"
	CHARSET   = "utf8"
)

// Init 数据库初始化
func Init() {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", USER_NAME, PASSWORD, HOST, PORT, DATABASE, CHARSET)
	database, err := sqlx.Open("mysql", dbDSN)
	if err != nil {
		fmt.Println("Open mysql fained!", err)
	}
	db = database
}

// GetConnection 获得数据库连接
func GetConnection() *sqlx.DB {
	return db
}
