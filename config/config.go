package config

import (
	"log"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct{
	DBHost string
	DBPort string 
	DBUser string
	DBPassword string 
	DBName string 

	RedisHost string 
	RedisPort string 
	RedisPassword string

	JWTSecret string 
	JWTExpiration int
	JWTSalt int

	UserIdSize int

	SMTPHOST string
	SMTPPORT string
	SMTPUSER string
	SMTPPASS string
	SMTPFORM string
}


func LoadConfig() *Config{
	err := godotenv.Load()
	if err != nil {
		log.Println("problem in loading env data")
	}
	
	JWT_SALT, _ := strconv.Atoi(os.Getenv("JWT_SALT"))
	UserIdSize, _ := strconv.Atoi(os.Getenv("UserIdSize"))

	epxirationHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if epxirationHours == 0 {epxirationHours= 24} // default value 
	return &Config{
        DBHost:     os.Getenv("DB_HOST"),
        DBPort:     os.Getenv("DB_PORT"),
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),

        RedisHost:     os.Getenv("REDIS_HOST"),
        RedisPort:     os.Getenv("REDIS_PORT"),
        RedisPassword: os.Getenv("REDIS_PASSWORD"),

        JWTSecret:     os.Getenv("JWT_SECRET"),
        JWTExpiration: epxirationHours,
		JWTSalt: JWT_SALT,

		UserIdSize: UserIdSize,


		SMTPHOST: os.Getenv("SMtpHost"),
		SMTPPORT: os.Getenv("smtpPort"),
		SMTPUSER: os.Getenv("smtpUser"),
		SMTPPASS: os.Getenv("smtpPass"),
		SMTPFORM: os.Getenv("smtpFrom"),
    }

}