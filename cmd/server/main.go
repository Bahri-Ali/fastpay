package main 
import(
	"fmt"
	"log"
	"fastpay-backend/config"
    "fastpay-backend/database"
)


func main(){
	ctx := config.LoadConfig()

	database.connectDb(ctx)
	database.connectRedis(ctx)

	fmt.Println("server are working")
}