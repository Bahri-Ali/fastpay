package database

import (
	"context"
	"fmt"
	"log"
	"fastpay-backend/config"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/redis/go-redis/v9"
)

var (
	pgPoll *pgxpool.Pool 
	rdb *redis.Client

)

func connectDb(ctx *config.Config){
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        ctx.DBHost, ctx.DBUser, ctx.DBPassword, ctx.DBName, ctx.DBPort)
	
	var err error
	pgPoll , err = pgxpool.New(context.Background() , dsn)
	if err != nil {log.Printf("db connection faild ")}

	err = pgPoll.Ping(context.Background())
	if err != nil {log.Println("problem of ping db")}

	log.Println("connect to db done ")
}

func connectRedis(ctx *config.Config){
	rdb = redis.NewClient(&redis.Options{
		Addr: ctx.RedisHost+":"+ctx.RedisPort,
		Password: ctx.RedisPassword,
		DB: 0,
	})
	_,err := rdb.Ping(context.Background()) .Result()
	if err != nil{log.Println("faild to connect to redis")}
	log.Println("connect to redis done")
}