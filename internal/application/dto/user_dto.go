package dto

import (
	"time"
	"yiwen/go-ddd/internal/domain/entity"
)

// DTO (Data Transfer ob) 数据传输对象
// DTO 用于层间数据传递，与领域实体分离
// 好处：
// 1. 隐藏领域模型的内部结构
// 2. 根据不同场景返回不同接口
// 3. 便于API版本演进

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Nickname string `json:"nickname" binding:"omitempty,min=3,max=50"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string  `json:"token"`
	ExpiresAt int64   `json:"expires_at"`
	User      UserDTO `json:"user"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=255"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UserDTO struct {
	ID       uint64    `json:"id"`
	UUID     string    `json:"uuid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Nickname string    `json:"nickname"`
	Avatar   string    `json:"avatar"`
	Status   int       `json:"status"`
	Role     string    `json:"role"`
	CreateAt time.Time `json:"create_at"`
}

type UserListDTO struct {
	Total int64     `json:"total"`
	Items []UserDTO `json:"items"`
}

func ToUserDTO(user *entity.User) UserDTO {
	return UserDTO{
		ID:       user.ID,
		UUID:     user.UUID,
		Username: user.Username,
		Email:    user.Email.String(),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Status:   int(user.Status),
		Role:     string(user.Role),
	}
}

func ToUserDTOList(users []*entity.User) []UserDTO {
	dtos := make([]UserDTO, len(users))
	for i, user := range users {
		dtos[i] = ToUserDTO(user)
	}
	return dtos
}

type PaginationRequest struct {
	Page     int `json:"page" binding:"min=1"`
	PageSize int `json:"page_size" binding:"min=1,max=100"`
}

func (p *PaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationRequest) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}
