package service

import (
	"context"
	"errors"
	"yiwen/go-ddd/internal/domain/entity"
	"yiwen/go-ddd/internal/domain/repository"
)

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserNotActive         = errors.New("user not active")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserNotAdmin          = errors.New("user is not an admin")
)

// UserDomainService 用户领域服务
// 领域服务用于处理不属于单个实体的业务逻辑
// 例如 涉及到多个实体的操作，需要访问仓储的验证逻辑
type UserDomainService struct {
	userRepo repository.UserRepository
}

// NewUserDomainService 创建用户领域服务
func NewUserDomainService(userRepo repository.UserRepository) *UserDomainService {
	return &UserDomainService{userRepo: userRepo}
}

// ValidateUniqueUsername 验证用户名是否唯一
func (s *UserDomainService) ValidateUniqueUsername(ctx context.Context, username string) error {
	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return err
	}
	if exists {
		return ErrUsernameAlreadyExists
	}
	return nil
}

func (s *UserDomainService) ValidateUniqueEmail(ctx context.Context, email string) error {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailAlreadyExists
	}
	return nil
}

func (s *UserDomainService) ValidateUserCredentials(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if err := user.Password.Verify(password); err != nil {
		return nil, ErrInvalidCredentials
	}
	if !user.IsActive() {
		return nil, ErrUserNotActive
	}
	return user, nil
}

func (s *UserDomainService) CanUserPerformAction(user *entity.User, action string) bool {
	if user.IsAdmin() {
		return true
	}

	alowedActions := map[string]bool{
		"view_profile":    true,
		"update_profile":  true,
		"change_password": true,
	}

	return alowedActions[action]
}

func (s *UserDomainService) TransferAdmin(ctx context.Context, userID uint64) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if !user.IsAdmin() {
		return ErrUserNotAdmin
	}
	user.PromoteToAdmin()
	return s.userRepo.Save(ctx, user)
}
