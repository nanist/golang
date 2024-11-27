package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func Getcode(len int) string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < len; i++ {
		num := rand.Intn(10)
		code += strconv.Itoa(num)
	}
	return code
}
