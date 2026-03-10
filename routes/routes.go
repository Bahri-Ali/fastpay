package routes
import (
    "fastpay-backend/internal/auth"
    "github.com/gin-gonic/gin"
)

type RouteConfig struct{
	AuthCntr *auth.Controller
}

func  SetupRouter(cfg *RouteConfig) *gin.Engine{
	Router := gin.Default()

	api := Router.Group("/api/v1")

    if cfg.AuthCntr !=nil{
    	AuthApi := api.Group("/auth")
		{
    	    AuthApi.POST("/register", cfg.AuthCntr.Register)
            AuthApi.POST("/login", cfg.AuthCntr.Login)
        }
    }


	return Router
}

