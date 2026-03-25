package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateVerificationCode 生成6位数字验证码
func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000 // 生成100000-999999之间的随机数
	return fmt.Sprintf("%06d", code)
}
