// package routes

// import (
// 	"fastpay-backend/internal/auth"
// 	"fastpay-backend/internal/middleware"


// 	// "fastpay-backend/internal/transaction"

// 	"github.com/gin-gonic/gin"
// )

// type RouteConfig struct{
// 	AuthCntr *auth.Controller
// 	AuthRepo *auth.Repository
 

// }

// func  SetupRouter(cfg *RouteConfig) *gin.Engine{
// 	Router := gin.Default()

// 	Router.Use(middleware.RateLimit())
// 	api := Router.Group("/api/v1")

//     if cfg.AuthCntr !=nil{
//     	AuthApi := api.Group("/auth")
// 		{
//     	    AuthApi.POST("/register", cfg.AuthCntr.Register)
//             AuthApi.POST("/login", cfg.AuthCntr.Login)
//         }

		
//     }
		
	

	
// 	return Router
// }

package routes

import (
    "fastpay-backend/internal/auth"
    "fastpay-backend/internal/middleware"
    "fastpay-backend/internal/transaction" // استيراد حزمة المعاملات الجديدة

    "github.com/gin-gonic/gin"
)

type RouteConfig struct {
    AuthCntr     *auth.Controller
    AuthRepo     auth.Repository      // تم التصحيح: استخدام Interface وليس Pointer
    TxController *transaction.Controller // إضافة متحكم المعاملات
}

func SetupRouter(cfg *RouteConfig) *gin.Engine {
    Router := gin.Default()

    // تطبيق Rate Limit على كامل التطبيق
    Router.Use(middleware.RateLimit())

    api := Router.Group("/api/v1")

    // 1. Public Routes (المسارات العامة)
    if cfg.AuthCntr != nil {
        AuthApi := api.Group("/auth")
        {
            AuthApi.POST("/register", cfg.AuthCntr.Register)
            AuthApi.POST("/login", cfg.AuthCntr.Login)
        }
    }

    // 2. Protected Routes (المسارات المحمية - تتطلب Token)
    if cfg.TxController != nil {
        protected := api.Group("/")
        // استخدام Middleware للتحقق من الـ Token
        protected.Use(middleware.AuthMiddleware(cfg.AuthRepo))
        {
            protected.POST("/transfer", cfg.TxController.InitTransfer)
            protected.POST("/transfer/verify", cfg.TxController.VerifyTransfer)
        }
    }

    return Router
}