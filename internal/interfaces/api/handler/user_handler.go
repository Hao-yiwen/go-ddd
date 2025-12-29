package handler

import "yiwen/go-ddd/internal/interfaces/api/middleware"

type UserHandler struct {
	userService *service.UserApplicationService
	jwtAuth     *middleware.JWTAuth
}

func NewUserHandler(userService *service.UserApplicationService, jwtAuth *middleware.JWTAuth) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtAuth:     jwtAuth,
	}
}
