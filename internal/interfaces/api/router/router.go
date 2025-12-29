package router

import (
	"net/http"
	"yiwen/go-ddd/internal/interfaces/api/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine      *gin.Engine
	userHandler *handler.userHandler
	jwtAuth     *middleware.JWTAuth
}

func NewRouter(userHandler *handler.userHandler, jwtAuth *middleware.JWTAuth) *Router {
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
