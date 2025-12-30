package query

// Query 查询模式
// CORS 中的查询部分 用于读操作
// 查询不改变系统状态，只返回数据

// GetUserByIDQuery 根据id查询用户
type GetUserByIDQuery struct {
	UserID uint64
}

// NewGetUserByIDQuery 创建根据id查询用户查询
func NewGetUserByIDQuery(userID uint64) *GetUserByIDQuery {
	return &GetUserByIDQuery{UserID: userID}
}

// GetUserByUUIDQuery 根据uuid查询用户
type GetUserByUUIDQuery struct {
	UUID string
}

// NewGetUserByUUIDQuery 创建根据uuid查询用户查询
func NewGetUserByUUIDQuery(uuid string) *GetUserByUUIDQuery {
	return &GetUserByUUIDQuery{UUID: uuid}
}

// GetUserByUsernameQuery 根据用户名查询用户
type GetUserByUsernameQuery struct {
	Username string
}

// NewGetUserByUsernameQuery 创建根据用户名查询用户查询
func NewGetUserByUsernameQuery(username string) *GetUserByUsernameQuery {
	return &GetUserByUsernameQuery{Username: username}
}

// ListUserQuery 查询用户列表
type ListUsersQuery struct {
	Offset int
	Limit  int
}

func NewListUserQuery(offset, limit int) *ListUsersQuery {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return &ListUsersQuery{Offset: offset, Limit: limit}
}

// LoginQuery 登录查询
type LoginQuery struct {
	Username string
	Password string
}

// NewLoginQuery 创建登陆查询
func NewLoginQuery(username, password string) *LoginQuery {
	return &LoginQuery{Username: username, Password: password}
}
