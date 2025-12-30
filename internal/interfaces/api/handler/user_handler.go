package handler

import (
	"net/http"
	"strconv"
	"yiwen/go-ddd/internal/application/command"
	"yiwen/go-ddd/internal/application/dto"
	"yiwen/go-ddd/internal/application/query"
	"yiwen/go-ddd/internal/application/service"
	"yiwen/go-ddd/internal/interfaces/api/middleware"

	"github.com/gin-gonic/gin"
)

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

// Register 用户注册
// POST /api/v1/users/register
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	cmd := command.NewRegisterUserCommand(req.Username, req.Email, req.Password, req.Nickname)
	user, err := h.userService.Register(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    200,
		"message": "User registered successfully",
		"data":    user,
	})
}

// Login 用户登录
// POST /api/v1/users/login
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	q := query.NewLoginQuery(req.Username, req.Password)
	user, err := h.userService.Login(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
		return
	}

	token, expiresAt, err := h.jwtAuth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Login successfully",
		"data": dto.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			User:      *user,
		},
	})
}

// GetUser 获取用户信息
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid user ID",
		})
		return
	}

	q := query.NewGetUserByIDQuery(id)
	user, err := h.userService.GetUserByID(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// ListUsers 获取用户列表
// GET /api/v1/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	q := query.NewListUserQuery(req.Page, req.PageSize)
	users, err := h.userService.ListUsers(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

// UpdateProfile 更新用户资料
// PUT /api/v1/users/:id
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid user ID",
		})
		return
	}

	currentUserID, _ := middleware.GetUserIDFromContext(c)
	if currentUserID != id {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "Forbidden",
		})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	cmd := command.NewUpdateProfileCommand(id, req.Nickname, req.Avatar)
	user, err := h.userService.UpdateProfile(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Profile updated successfully",
		"data":    user,
	})
}

// ChangePassword 修改密码
// PUT /api/v1/users/:id/change-password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid user ID",
		})
		return
	}

	currentUserID, _ := middleware.GetUserIDFromContext(c)
	if currentUserID != id {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "Forbidden",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	cmd := command.NewChangePasswordCommand(id, req.OldPassword, req.NewPassword)
	err = h.userService.ChangePassword(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Password changed successfully",
	})
}

// DeleteUser 删除用户
// DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid user ID",
		})
		return
	}

	cmd := command.NewDeleteUserCommand(id)
	if err := h.userService.DeleteUser(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "User deleted successfully",
	})
}

// GetCurrentUser 获取当前用户信息
// GET /api/v1/users/me
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Unauthorized",
		})
		return
	}

	q := query.NewGetUserByIDQuery(userID)
	user, err := h.userService.GetUserByID(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Current user retrieved successfully",
		"data":    user,
	})

}
