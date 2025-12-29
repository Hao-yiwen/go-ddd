package model

import (
	"time"
	"yiwen/go-ddd/internal/domain/entity"
	"yiwen/go-ddd/internal/domain/valueobject"

	"gorm.io/gorm"
)

// UserModel 用户数据库模型
// 数据库模型与领域实体分离的好处：
// 1. 领域实体不受数据库结构影响
// 2. 可以自由添加数据库特有字段 - 例如软删除
// 3. 便于处理ORM特有的变迁和钩子
type UserModel struct {
	ID           uint64    `gorm:"primayKey;autoIncrement"`
	UUID         string    `gorm:"type:varchar(36);uniqueIndex;not null"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Nickname     string    `gorm:"type:varchar(50)"`
	Avatar       string    `gorm:"type:varchar(255)"`
	Status       int       `gorm:"type:tinyint(1);not null;default:1"` // tinyint(1) 是 MySQL 的字段类型，适合用于布尔值或者较小的整型状态字段
	Role         string    `gorm:"type:varchar(20);not null;default:user"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	// 软删除字段，gorm内置类型，表示删除时间。被删除不会真正移除，只是设置删除时间。
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string {
	return "users"
}

func (m *UserModel) ToEnitity() *entity.User {
	email, _ := valueobject.NewEmail(m.Email)
	password := valueobject.NewPasswordFromHash(m.PasswordHash)

	return &entity.User{
		ID:       m.ID,
		UUID:     m.UUID,
		Username: m.Username,
		Email:    email,
		Password: password,
		Nickname: m.Nickname,
		Avatar:   m.Avatar,
		Status:   entity.UserStatus(m.Status),
		Role:     entity.UserRole(m.Role),
	}
}

func (m *UserModel) FromEntity(user *entity.User) {
	m.ID = user.ID
	m.UUID = user.UUID
	m.Username = user.Username
	m.Email = user.Email.String()
	m.PasswordHash = user.Password.Hash()
	m.Nickname = user.Nickname
	m.Avatar = user.Avatar
	m.Status = int(user.Status)
	m.Role = string(user.Role)
}
