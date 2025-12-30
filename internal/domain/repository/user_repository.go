package repository

import (
	"context"
	"yiwen/go-ddd/internal/domain/entity"
)

// UserRepository 用户仓库接口
// 仓储模式是DDD中的重要模式：
// 1. 领域层只定义接口，不关心具体实现
// 2. 基础设施层提供具体实现 如mysql redis
// 3. 这样可以实现依赖倒置， 领域层不依赖具体技术
// 4. 便于测试和扩展
type UserRepository interface {
	// save 保存用户
	Save(ctx context.Context, user *entity.User) error

	// FindByID 根据id查询用户
	FindByID(ctx context.Context, id uint64) (*entity.User, error)

	// FindByUUID 根据uuid查找用户
	FindByUUID(ctx context.Context, uuid string) (*entity.User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// Delete 删除用户 软删除
	Delete(ctx context.Context, id uint64) error

	// List 分页查询用户列表
	List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)

	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
