package main

import (
	"fastpay-backend/config"
	"fastpay-backend/database"
	"fastpay-backend/internal/auth"
	"fastpay-backend/internal/transaction"
	"fastpay-backend/internal/user"
	"fastpay-backend/internal/wallet"
	email "fastpay-backend/pkg/emails"
	"fastpay-backend/routes"
	"fmt"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("DEBUG CONFIG: Host=%s, Port=%s\n", cfg.SMTPHost, cfg.SMTPPort)
	database.ConnectDb(cfg)
	database.ConnectRedis(cfg)

	walletRepo := wallet.NewRepository(database.PgPoll)
	authRepo := auth.NewRepository(database.PgPoll)
	txRepo := transaction.NewRepository(database.PgPoll)

	mailer := email.NewMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)

	walletService := wallet.NewService(walletRepo)
	authService := auth.NewService(authRepo, walletRepo, cfg)
	txService := transaction.NewService(txRepo, database.Rdb, mailer)

	walletController := wallet.NewController(walletService)
	authController := auth.NewController(authService)
	txController := transaction.NewController(txService)

	userRepo := user.NewRepository(database.PgPoll)
	userService := user.NewService(userRepo, database.Rdb, mailer)
	userController := user.NewController(userService)

	routerConfig := &routes.RouteConfig{
		AuthCntr:         authController,
		AuthRepo:         authRepo,
		TxController:     txController,
		UserController:   userController,
		WalletController: walletController,
	}

	router := routes.SetupRouter(routerConfig)

	log.Println("Server starting on :8080 with SSL (HTTPS)")

	if err := router.RunTLS(":8080", "cert.pem", "key.pem"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
