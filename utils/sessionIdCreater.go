package utils

import (
	"fmt"
	"github.com/satori/go.uuid"
)

func GetNewSessionId() string {
	ul := uuid.NewV4()
	fmt.Println("生成SessionId：", ul)
	return ul.String()
}
