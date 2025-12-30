package mysql

import (
	"context"
	"errors"
	"yiwen/go-ddd/internal/domain/entity"
	"yiwen/go-ddd/internal/domain/repository"
	"yiwen/go-ddd/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

// UserRepository Mysql用户仓库实现
// 这里是仓库接口具体实现
// 基础设施实现领域层定义的接口
type UserRepository struct {
	db *gorm.DB
}

// FindByUsername implements [repository.UserRepository].
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	panic("unimplemented")
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

// Save 方法将领域层的 User 实体保存到数据库。
// 首先，利用 model.FromEntity 将领域实体转换为数据库模型（UserModel），这样确保领域层与基础设施层解耦。
// 如果 user.ID == 0，说明这是一个新用户，还没有主键ID（ID通常由数据库自增生成），
// 因此调用 Create 方法插入新记录，并插入后将数据库生成的 ID 回填到实体的 user.ID 字段。
// 如果 user.ID != 0，说明该用户已存在，这是一次更新操作，调用 Save 方法。
// 这样即可根据 user 是否有 ID 自动区分新增和更新。
// 若数据库操作出错，直接返回错误。
// 一切操作都通过 GORM 的 WithContext 保证支持 trace、timeout、cancel 等。
func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
	userModel := model.FromEntity(user)

	if user.ID == 0 {
		// 新建用户，插入数据库
		if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
			return err
		}
		user.ID = userModel.ID // 回写自增ID到实体
	} else {
		// 已有用户，更新数据库
		if err := r.db.WithContext(ctx).Save(userModel).Error; err != nil {
			return err
		}
	}
	return nil
}

// FindByID 方法根据数据库ID查询用户。
func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	var userModel model.UserModel

	if err := r.db.WithContext(ctx).First(&userModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return userModel.ToEnitity(), nil
}

func (r *UserRepository) FindByUUID(ctx context.Context, uuid string) (*entity.User, error) {
	var userModel model.UserModel

	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return userModel.ToEnitity(), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel model.UserModel

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return userModel.ToEnitity(), nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.UserModel{}, id).Error
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	var userModels []model.UserModel
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.UserModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("id DESC").
		Find(&userModels).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*entity.User, len(userModels))
	for i, model := range userModels {
		users[i] = model.ToEnitity()
	}

	return users, total, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
