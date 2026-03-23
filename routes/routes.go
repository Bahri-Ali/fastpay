package routes

import (
	"fastpay-backend/database"
	"fastpay-backend/internal/auth"
	"fastpay-backend/internal/middleware"
	"fastpay-backend/internal/transaction"
	"fastpay-backend/internal/user"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RouteConfig struct {
    AuthCntr     *auth.Controller
    AuthRepo     auth.Repository      
    TxController *transaction.Controller 
    UserController *user.Controller
    Rdb *redis.Client
}

func SetupRouter(cfg *RouteConfig) *gin.Engine {
    Router := gin.Default()

    
    Router.Use(middleware.RateLimit())

    api := Router.Group("/api/v1")

    if cfg.AuthCntr != nil {
        AuthApi := api.Group("/auth")
        {
            AuthApi.POST("/register", cfg.AuthCntr.Register)
            AuthApi.POST("/login", cfg.AuthCntr.Login)
        }
    }
    if database.Rdb == nil {
    log.Println("❌ Redis is NIL")
}
    if cfg.TxController != nil {
        protected := api.Group("/")
        protected.Use(middleware.AuthMiddleware(cfg.AuthRepo))
        {
            protected.POST("/transfer", cfg.TxController.InitTransfer)
            protected.POST("/transfer/verify", cfg.TxController.VerifyTransfer)
            protected.GET("/GetTransactions", cfg.TxController.GetTransactionHistory)
            protected.GET("/ws/transactions", func(c *gin.Context) {
            transaction.HandleWebSocket(c,database.Rdb) 
})
        }
    }

    if cfg.UserController != nil{
        protected := api.Group("/user")
        protected.Use(middleware.AuthMiddleware(cfg.AuthRepo))
        {
            protected.GET("/profile" , cfg.UserController.GetProfile )
            protected.POST("/change-password/init" , cfg.UserController.ChangePasswordInit )
            protected.POST("/change-password/verify" , cfg.UserController.ChangePasswordVerify )
        }
    }


    return Router
}