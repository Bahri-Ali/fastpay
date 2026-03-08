package utils

import (
	"fastpay-backend/config"
	"math/rand"

)

func GenerateUserID(wilayaCode int , birthDate int ,ctx *config.Config) string{
	const digits = "0123456789"
    b := make([]byte, ctx.UserIdSize)

    // dateStr := birthDate.Format("20060102")

    for i := range b {
        b[i] = digits[rand.Intn(len(digits))]
    }
    randomStr := string(b)
	return string(wilayaCode) + randomStr
}