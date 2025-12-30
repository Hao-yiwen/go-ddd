package router

import (
	"net/http"
	"yiwen/go-ddd/internal/interfaces/api/handler"
	"yiwen/go-ddd/internal/interfaces/api/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine      *gin.Engine
	userHandler *handler.UserHandler
	jwtAuth     *middleware.JWTAuth
}

func NewRouter(userHandler *handler.UserHandler, jwtAuth *middleware.JWTAuth) *Router {
	return &Router{
		engine:      gin.New(),
		userHandler: userHandler,
		jwtAuth:     jwtAuth,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(gin.Logger())
	r.engine.Use(gin.Recovery())
	r.engine.Use(CORSMiddleware())

	r.engine.GET("/heath", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "OK",
		})
	})

	v1 := r.engine.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/register", r.userHandler.Register)
			users.POST("/login", r.userHandler.Login)

			authUsers := users.Group("")
			authUsers.Use(r.jwtAuth.AdminMiddleware())
			{
				authUsers.GET("/:id", r.userHandler.GetUser)
				authUsers.PUT("/:id", r.userHandler.UpdateProfile)
				authUsers.PUT("/:id/change-password", r.userHandler.ChangePassword)
				authUsers.DELETE("/:id", r.userHandler.DeleteUser)
				authUsers.GET("/me", r.userHandler.GetCurrentUser)
			}

			//管理员接口
			adminUsers := users.Group("")
			adminUsers.Use(r.jwtAuth.AuthMiddleware())
			adminUsers.Use(r.jwtAuth.AdminMiddleware())
			{
				adminUsers.GET("/", r.userHandler.ListUsers)
				adminUsers.DELETE("/:id", r.userHandler.DeleteUser)
			}
		}
	}

	return r.engine
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
