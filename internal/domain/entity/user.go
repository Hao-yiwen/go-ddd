package entity

import (
	"time"
	"yiwen/go-ddd/internal/domain/valueobject"
)

type UserStatus int

const (
	UserStatusActive   UserStatus = 1 // 激活
	UserStatusInactive UserStatus = 2 // 未激活
	UserStatusBanned   UserStatus = 3 // 禁用
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// user 用户实体
// 实体是ddd中的核心概念 具有唯一标识 id
// 实体的相等性由id决定 而不是属性
type User struct {
	ID        uint64               // 数据库自增ID
	UUID      string               // 业务唯一标识
	Username  string               // 用户名
	Email     valueobject.Email    // 邮箱
	Password  valueobject.Password //密码
	Nickname  string               // 昵称
	Avatar    string               // 头像
	Status    UserStatus           // 状态
	Role      UserRole             // 角色
	CreatedAt time.Time            // 创建时间
	UpdatedAt time.Time            // 更新时间
	DeletedAt time.Time            // 删除时间
}

func NewUser(uuid string, username string, email valueobject.Email, password valueobject.Password) *User {
	return &User{
		UUID:      uuid,
		Username:  username,
		Email:     email,
		Password:  password,
		Status:    UserStatusActive,
		Role:      UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u *User) UpdateProfile(nickname, avatar string) {
	u.Nickname = nickname
	u.Avatar = avatar
	u.UpdatedAt = time.Now()
}

func (u *User) ChangePassword(newPassword valueobject.Password) {
	u.Password = newPassword
	u.UpdatedAt = time.Now()
}

func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
}

func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now()
}

func (u *User) Ban() {
	u.Status = UserStatusBanned
	u.UpdatedAt = time.Now()
}

func (u *User) PromoteToAdmin() {
	u.Role = UserRoleAdmin
	u.UpdatedAt = time.Now()
}
