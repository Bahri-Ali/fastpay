package routes

import (
	"fastpay-backend/internal/auth"
	"fastpay-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct{
	AuthCntr *auth.Controller
	AuthRepo *auth.Repository
}

func  SetupRouter(cfg *RouteConfig) *gin.Engine{
	Router := gin.Default()

	Router.Use(middleware.RateLimit())
	api := Router.Group("/api/v1")

    if cfg.AuthCntr !=nil{
    	AuthApi := api.Group("/auth")
		{
    	    AuthApi.POST("/register", cfg.AuthCntr.Register)
            AuthApi.POST("/login", cfg.AuthCntr.Login)
        }
    }

	protected:= api.Group("/")
	protected.Use(middleware.AuthMiddleware(*cfg.AuthRepo))
	{
		 // protected.GET("/profile", cfg.AuthCntr.GetProfile)
	}
	return Router
}

