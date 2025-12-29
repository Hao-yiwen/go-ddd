package mysql

import "gorm.io/gorm"

// UserRepository Mysql用户仓库实现
// 这里是仓库接口具体实现
// 基础设施实现领域层定义的接口
type UserRepository struct {
	db *gorm.DB
}
