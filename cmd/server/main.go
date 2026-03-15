package main

import (
	"fastpay-backend/config"
	"fastpay-backend/database"
	"fastpay-backend/internal/auth"
	"fastpay-backend/internal/transaction" 
	"fastpay-backend/internal/wallet"
	"fastpay-backend/pkg/emails"
	"fastpay-backend/routes"
	"log"
)

func main() {
    // 1. Load Config
    cfg := config.LoadConfig()

    // 2. Connect Databases
    database.ConnectDb(cfg)
    database.ConnectRedis(cfg)

    // 3. Initialize Repositories
    walletRepo := wallet.NewRepository(database.PgPoll)
    authRepo := auth.NewRepository(database.PgPoll)
    txRepo := transaction.NewRepository(database.PgPoll) // جديد

    // 4. Initialize Utils (Mailer)
    // تأكد من أن الـ Config يحتوي على بيانات SMTP
    mailer := email.NewMailer(cfg.SMTPHOST, cfg.SMTPPORT, cfg.SMTPUSER, cfg.SMTPPASS, cfg.SMTPFORM)

    // 5. Initialize Services
    authService := auth.NewService(authRepo, walletRepo, cfg)
    txService := transaction.NewService(txRepo, database.Rdb, mailer) // جديد: تمرير Redis و Mailer

    // 6. Initialize Controllers
    authController := auth.NewController(authService)
    txController := transaction.NewController(txService) // جديد

    // 7. Setup Router Config
    routerConfig := &routes.RouteConfig{
        AuthCntr:     authController,
        AuthRepo:     authRepo,    // تم التصحيح: تمرير الـ Interface مباشرة
        TxController: txController, // جديد
    }

    router := routes.SetupRouter(routerConfig)

    // 8. Run Server
    log.Println("Server starting on :8080 with SSL (HTTPS)")
    
    // تصحيح خطأ الـ Syntax في الـ if
    if err := router.RunTLS(":8080", "cert.pem", "key.pem"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}