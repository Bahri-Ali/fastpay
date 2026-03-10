package main

import (
	"fastpay-backend/config"
	"fastpay-backend/database"
	"fastpay-backend/internal/auth"
	"fastpay-backend/routes"
	"log"
)


func main(){
	cfg := config.LoadConfig()

	database.ConnectDb(cfg)
	database.ConnectRedis(cfg)

	//auth layer
	authRepo := auth.NewRepository(database.PgPoll)
	authService := auth.NewService(authRepo  , cfg)
	authController := auth.NewController(authService)



	routerConfig := &routes.RouteConfig{
		AuthCntr: authController,
	}	
	

	router := routes.SetupRouter(routerConfig)
    log.Println("Server starting on :8080")
        if err := router.Run(":8080"); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
}