package hash

import (
	"fastpay-backend/config"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(ctx *config.Config , password string)(string, error){
	bytes , err := bcrypt.GenerateFromPassword([]byte(password) , ctx.JWTSalt)
	if err != nil{fmt.Println("err in password hash function ")}
	return string(bytes) , nil
}

func CheckPasswordHash(password string , hashPassword string)bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword) , []byte(password))
	return err == nil 
}