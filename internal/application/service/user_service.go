package service

import (
	"context"
	"yiwen/go-ddd/internal/application/command"
	"yiwen/go-ddd/internal/application/dto"
	"yiwen/go-ddd/internal/application/query"
	"yiwen/go-ddd/internal/domain/aggregate"
	"yiwen/go-ddd/internal/domain/repository"
	"yiwen/go-ddd/internal/domain/service"
	"yiwen/go-ddd/internal/domain/valueobject"
	"yiwen/go-ddd/pkg/errors"

	domainservice "yiwen/go-ddd/internal/domain/service"

	"github.com/google/uuid"
)

// UserApplicationService
// 应用服务是应用层的核心 负责
// 1. 协调领域层和基础设施层
// 2. 处理事物
// 3. 调用领域服务
// 4. 不包含业务逻辑
type UserApplicationService struct {
	userRepo          repository.UserRepository
	userDomainService service.UserDomainService
}

// NewUserApplicationService 创建用户应用服务
func NewUserApplicationService(userRepo repository.UserRepository, userDomainService domainservice.UserDomainService) *UserApplicationService {
	return &UserApplicationService{userRepo: userRepo, userDomainService: userDomainService}
}

// Register 注册用户
func (s *UserApplicationService) Register(ctx context.Context, cmd *command.RegisterUserCommand) (*dto.UserDTO, error) {
	// 验证用户名是否唯一
	if err := s.userDomainService.ValidateUniqueUsername(ctx, cmd.Username); err != nil {
		return nil, err
	}

	// 验证邮箱是否唯一
	if err := s.userDomainService.ValidateUniqueEmail(ctx, cmd.Email); err != nil {
		return nil, err
	}

	email, err := valueobject.NewEmail(cmd.Email)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid email")
	}

	password, err := valueobject.NewPassword(cmd.Password)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid password")
	}

	userAggregate := aggregate.Register(uuid.New().String(), cmd.Username, email, password)
	userAggregate.User.Nickname = cmd.Nickname

	if err := s.userRepo.Save(ctx, userAggregate.User); err != nil {
		return nil, errors.Wrapf(err, "failed to save user")
	}

	// 发布领域事件

	result := dto.ToUserDTO(userAggregate.User)
	return &result, nil
}

func (s *UserApplicationService) Login(ctx context.Context, q *query.LoginQuery) (*dto.UserDTO, error) {
	user, err := s.userDomainService.ValidateUserCredentials(ctx, q.Username, q.Password)
	if err != nil {
		return nil, err
	}

	result := dto.ToUserDTO(user)
	return &result, nil
}

func (s *UserApplicationService) GetUserByID(ctx context.Context, q *query.GetUserByIDQuery) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, q.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "user not found")
	}

	result := dto.ToUserDTO(user)
	return &result, nil
}

func (s *UserApplicationService) ListUsers(ctx context.Context, q *query.ListUsersQuery) (*dto.UserListDTO, error) {
	users, total, err := s.userRepo.List(ctx, q.Offset, q.Limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list users")
	}

	dtos := dto.ToUserDTOList(users)
	return &dto.UserListDTO{
		Total: total,
		Items: dtos,
	}, nil
}

func (s *UserApplicationService) UpdateProfile(ctx context.Context, cmd *command.UpdateProfileCommand) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "user not found")
	}

	userAggregate := aggregate.NewUserAggregate(user)
	userAggregate.UpdateProfile(cmd.Nickname, cmd.Avatar)

	if err := s.userRepo.Save(ctx, userAggregate.User); err != nil {
		return nil, errors.Wrap(err, "failed to save user")
	}

	result := dto.ToUserDTO(userAggregate.User)
	return &result, nil
}

func (s *UserApplicationService) ChangePassword(ctx context.Context, cmd *command.ChangePasswordCommand) error {
	user, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return errors.Wrap(err, "user not found")
	}

	if err := user.Password.Verify(cmd.OldPassword); err != nil {
		return domainservice.ErrInvalidCredentials
	}

	newPassword, err := valueobject.NewPassword(cmd.NewPassword)
	if err != nil {
		return errors.Wrapf(err, "invalid new password")
	}

	userAggregate := aggregate.NewUserAggregate(user)
	userAggregate.ChangePassword(newPassword)

	if err := s.userRepo.Save(ctx, userAggregate.User); err != nil {
		return errors.Wrap(err, "failed to save user")
	}

	return nil
}

func (s *UserApplicationService) DeleteUser(ctx context.Context, cmd *command.DeleteUserCommand) error {
	if err := s.userRepo.Delete(ctx, cmd.UserID); err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}
