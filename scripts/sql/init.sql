-- =====================================================
-- DDD 用户管理系统 - 数据库初始化脚本
-- =====================================================

--- 创建数据库
CREATE DATABASE IF NOT EXISTS go_ddd DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

--- ==============================
--- 用户表
--- ==============================
CREAT TABLE IF NOT EXISTS users (
    --- 主键: 数据库自增ID
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',

    --- 业务唯一标识: UUID
    uuid VARCHAR(36) NOT NULL COMMENT '业务唯一标识UUID',

    --- 用户基本信息
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    email VARCHAR(100) NOT NULL COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希值',
    nickname VARCHAR(50) DEFAULT '' COMMENT '昵称',
    avatar VARCHAR(255) DEFAULT '' COMMENT '头像URL',

    --- 状态 和 角色
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1-激活 2-未激活 3-禁用',
    role VARCHAR(20) NOT NULL DEFAULT 'user' COMMENT '角色: user-普通用户 admin-管理员',

    --- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间（软删除）',

    --- 唯一索引
    UNIQUE INDEX uk_uuid(uuid),
    UNIQUE INDEX uk_username(username),
    UNIQUE INDEX uk_email(email),

    -- 普通索引
    INDEX idx_status(status),
    INDEX idx_role(role),
    INDEX idx_created_at(created_at),
    INDEX idx_deleted_at(deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT='用户表';

--- ==============================
--- 插入测试管理员账户
--- 密码: Admin123 (bcrypt加密)
--- ==============================
INSERT INTO users (uuid, username, email, password_hash, nickname, status, role)
VALUES (
    UUID(),
    'admin',
    'admin@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'Adminstrator',
    1,
    'admin'
) ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP;

--- ==============================
--- 说明
--- ==============================
--- 
--- 1. 表设计遵循DDD原则：
---    - id: 数据库层面的标识
---    - uuid: 业务层面的唯一标识 (领域层使用)
---    - password_hash: 存储加密后的密码，不存明文
---
--- 2. 软删除机制:
---    - deleted_at: 基于软删除
---    - grom会自动处理软删除逻辑，无需手动处理
--- 3. 状态说明:
---    - 1: 激活 (active) - 正常使用
---    - 2: 未激活 (inactive) - 需要激活
---    - 3: 禁用 (banned) - 被管理员禁用
---
--- 4. 角色说明
---   - user: 普通用户
---   - admin: 管理员
--- 5. 测试账户
---    - 用户名: admin
---    - 密码: Admin123
---    - 生产环境需修改或删除该账户