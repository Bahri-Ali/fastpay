package utils

import (
	"fastpay-backend/config"
	"fmt"
	"math/rand"
	"time"
)

func GenerateUserID(wilayaCode int , ctx *config.Config) string{
	const digits = "0123456789"
    b := make([]byte, ctx.UserIdSize)


    for i := range b {
        b[i] = digits[rand.Intn(len(digits))]
    }
    randomStr := string(b)
	return string(wilayaCode) + randomStr
}

func GeneratePIN() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%04d", rand.Intn(10000))
}