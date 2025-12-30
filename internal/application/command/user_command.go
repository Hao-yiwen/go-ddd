package command

// Command 命令模式
// CORS 命令查询职责分离
// 命令用于写操作，改变系统状态
// 命令通常不返回数据 或返回ID

// RegisterUserCommand 注册用户命令
type RegisterUserCommand struct {
	Username string
	Email    string
	Password string
	Nickname string
}

// NewRegisterUserCommand 创建注册用户命令
func NewRegisterUserCommand(username, email, password, nickname string) *RegisterUserCommand {
	return &RegisterUserCommand{Username: username, Email: email, Password: password, Nickname: nickname}
}

// UpdateProfileCommand 更新资料命令
type UpdateProfileCommand struct {
	UserID   uint64
	Nickname string
	Avatar   string
}

// NewUpdateProfileCommand 创建更新资料命令
func NewUpdateProfileCommand(userID uint64, nickname, avatar string) *UpdateProfileCommand {
	return &UpdateProfileCommand{UserID: userID, Nickname: nickname, Avatar: avatar}
}

// ChangePasswordCommand 修改密码命令
type ChangePasswordCommand struct {
	UserID      uint64
	OldPassword string
	NewPassword string
}

// NewChangePasswordCommand 创建修改密码命令
func NewChangePasswordCommand(userID uint64, oldPassword, newPassword string) *ChangePasswordCommand {
	return &ChangePasswordCommand{UserID: userID, OldPassword: oldPassword, NewPassword: newPassword}
}

// DeleteUserCommand 删除用户命令
type DeleteUserCommand struct {
	UserID uint64
}

// NewDeleteUserCommand 创建删除用户命令
func NewDeleteUserCommand(userID uint64) *DeleteUserCommand {
	return &DeleteUserCommand{UserID: userID}
}

// BanUserCommand 禁用用户命令
type BanUserCommand struct {
	UserID uint64
}

// NewBanUserCommand 创建禁用用户命令
func NewBanUserCommand(userID uint64) *BanUserCommand {
	return &BanUserCommand{UserID: userID}
}

// PromoteToAdminCommand 提升为管理员命令
type PromoteToAdminCommand struct {
	UserID uint64
}

// NewPromoteToAdminCommand 创建提升为管理员命令
func NewPromoteToAdminCommand(userID uint64) *PromoteToAdminCommand {
	return &PromoteToAdminCommand{UserID: userID}
}
