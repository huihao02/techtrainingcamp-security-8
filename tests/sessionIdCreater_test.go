package tests

import (
	"example/utils"
	"fmt"
	"testing"
)

func TestGetNewSessionId(t *testing.T) {
	fmt.Println(utils.GetNewSessionId())
}
