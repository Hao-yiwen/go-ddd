# go-ddd

用 Go 写的用户管理 HTTP API，按领域驱动设计（DDD）分层，并用命令 / 查询把写操作和读操作分开。这是一份学习记录：注释里保留了各层职责的说明，接口和领域对象比当前能跑通的路径更全。

[![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-1.11.0-008ECF)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/GORM-1.31.1-orange)](https://gorm.io/)

模块路径：`yiwen/go-ddd`。

## 这个仓库里有什么

已实现并接到 HTTP 的能力：

- 用户注册、登录（登录成功后签发 JWT）
- 按数字 ID 查询用户、查询当前用户、分页列出用户
- 更新昵称和头像、修改密码、软删除用户
- 邮箱、密码强度校验（bcrypt 保存哈希）
- MySQL 仓储（GORM），以及一份建表脚本

领域层里还有激活、停用、封禁、提升为管理员等领域方法和对应事件类型，应用层也声明了 `BanUserCommand`、`PromoteToAdminCommand` 以及按 UUID / 用户名查询的 Query。这些类型没有接到路由或应用服务，进程内也不会发布领域事件。

## 技术栈

| 用途 | 依赖 | 版本 |
| --- | --- | --- |
| 语言 | Go | 1.25.5（`go.mod`） |
| HTTP | `github.com/gin-gonic/gin` | v1.11.0 |
| ORM | `gorm.io/gorm` | v1.31.1 |
| MySQL 驱动 | `gorm.io/driver/mysql` | v1.6.0 |
| 配置 | `github.com/spf13/viper` | v1.21.0 |
| JWT | `github.com/golang-jwt/jwt/v5` | v5.3.0 |
| 密码哈希 | `golang.org/x/crypto`（bcrypt） | v0.40.0 |
| UUID | `github.com/google/uuid` | v1.6.0 |

仓库中没有测试文件、Dockerfile、Makefile 或前端页面。

## 目录结构

```text
.
├── cmd/api/main.go                         # 进程入口：读配置、连库、组装依赖、启动 HTTP
├── config/config.yaml                      # 示例配置（端口、MySQL、JWT）
├── internal
│   ├── domain                              # 领域层，不依赖 Gin / GORM
│   │   ├── entity/                         # User 实体、状态与角色
│   │   ├── aggregate/                      # 用户聚合根，收集领域事件
│   │   ├── valueobject/                    # Email、Password
│   │   ├── event/                          # 领域事件与发布接口（无实现）
│   │   ├── repository/                     # UserRepository 接口
│   │   └── service/                        # 用户名/邮箱唯一、凭证校验
│   ├── application                         # 用例编排，不含 HTTP 细节
│   │   ├── command/                        # 写操作命令
│   │   ├── query/                          # 读操作查询
│   │   ├── dto/                            # 请求体与响应用 DTO
│   │   └── service/                        # UserApplicationService
│   ├── infrastructure
│   │   ├── config/                         # Viper 加载配置
│   │   └── persistence
│   │       ├── model/                      # users 表的 GORM 模型
│   │       └── mysql/                      # UserRepository 的 MySQL 实现
│   └── interfaces/api
│       ├── handler/                        # HTTP 处理函数
│       ├── middleware/                     # JWT 解析、管理员检查、CORS
│       └── router/                         # 路由注册
├── pkg/errors/                             # 通用错误包装与业务错误码
├── scripts/sql/init.sql                    # 建库、建表、插入管理员的脚本
├── go.mod
└── go.sum
```

请求从 `interfaces` 进入，转成 command 或 query，由应用服务调用领域服务和仓储。仓储接口定义在领域层，MySQL 实现放在基础设施层。

## 领域规则

注册时会做这些检查（通过后用户状态为激活 `1`，角色为 `user`）：

- 用户名、邮箱在库中必须唯一
- 邮箱去掉首尾空白并转为小写，需匹配 `name@domain.tld` 形式
- 密码至少 8 位，且同时包含大写字母、小写字母和数字；入库的是 bcrypt 哈希（`bcrypt.DefaultCost`）

状态与角色取值：

| 字段 | 值 |
| --- | --- |
| `status` | `1` 激活，`2` 未激活，`3` 禁用 |
| `role` | `user` 普通用户，`admin` 管理员 |

登录会按用户名查用户、校验密码，并拒绝非激活用户。聚合上的注册、改资料、改密码会往内存里的事件切片追加事件，没有事件总线或订阅者把它们发出去。

## HTTP 接口

默认监听 `:8080`，统一 JSON 形如 `{"code":..., "message":..., "data":...}`。参数不合法时返回 HTTP 400、`Invalid request`（或 `Invalid user ID`）。应用服务返回的错误目前都映射成 HTTP 500、`Internal server error`，领域错误原文不会出现在响应里。

跨域：允许任意 Origin，方法为 `GET, POST, PUT, DELETE, OPTIONS`，请求头允许 `Content-Type` 和 `Authorization`。

| 方法 | 路径 | 中间件（按注册顺序） | 说明 |
| --- | --- | --- | --- |
| `GET` | `/heath` | 无 | 存活检查，返回 `code=200`、`message=OK` |
| `POST` | `/api/v1/users/register` | 无 | 注册。成功时 HTTP 状态为 201，响应体里的 `code` 仍是 200 |
| `POST` | `/api/v1/users/login` | 无 | 登录，返回 `token`、`expires_at` 和用户信息 |
| `GET` | `/api/v1/users/:id` | 仅 `AdminMiddleware` | 按数字 ID 查询 |
| `PUT` | `/api/v1/users/:id` | 仅 `AdminMiddleware` | 更新昵称、头像；处理函数还会核对 token 中的用户 ID |
| `PUT` | `/api/v1/users/:id/change-password` | 仅 `AdminMiddleware` | 修改密码，同样核对用户 ID |
| `DELETE` | `/api/v1/users/:id` | 仅 `AdminMiddleware`，以及下面管理员组再注册一次 | 软删除 |
| `GET` | `/api/v1/users/me` | 仅 `AdminMiddleware` | 当前用户 |
| `GET` | `/api/v1/users/` | `AuthMiddleware` 然后 `AdminMiddleware` | 用户列表 |

需要登录的接口使用请求头 `Authorization: Bearer <token>`。JWT 为 HS256，声明里包含 `user_id`、`username`、`role`、过期时间和签发者。

`AdminMiddleware` 只读取上下文里的 `role`，它本身不解析 token。上表里只挂了 `AdminMiddleware` 的路由，在进入处理函数之前就会得到 HTTP 403。列表和其中一条删除路由先经过 `AuthMiddleware`。

请求体：

```json
// POST /api/v1/users/register
{ "username": "yiwen", "email": "yiwen@example.com", "password": "Admin123", "nickname": "一文" }

// POST /api/v1/users/login
{ "username": "yiwen", "password": "Admin123" }

// PUT /api/v1/users/:id
{ "nickname": "新昵称", "avatar": "https://example.com/a.png" }

// PUT /api/v1/users/:id/change-password
{ "old_password": "Admin123", "new_password": "Admin456" }
```

校验标签：`username` 3–50；`email` 为邮箱；注册和修改密码时 `password` / `new_password` 至少 8 位（强度规则在领域层）；`nickname` 注册时若填写则为 3–50，更新时最长 50；`avatar` 最长 255。

用户 JSON 字段：`id`、`uuid`、`username`、`email`、`nickname`、`avatar`、`status`、`role`、`create_at`。`ToUserDTO` 没有给 `create_at` 赋值，序列化结果是时间零值。

列表的查询参数结构体使用的是 `json` 标签，`ShouldBindQuery` 读的是 query/form。处理函数又把 `page`、`page_size` 直接传给把它们当作 offset、limit 的查询对象。分页辅助方法 `GetOffset` / `GetLimit` 没有被调用。

## 配置

`config/config.yaml`：

```yaml
app:
  name: go-ddd
  port: 8080
  mode: debug          # 同时用作 GORM 日志级别：debug=Info，test=Warn，production=Error

database:
  host: localhost
  port: 3306
  username: root
  password: root
  database: go_ddd
  max_idle_conns: 10
  max_open_conns: 100

jwt:
  secret: your-super-secret-key-change-in-production
  expire_hour: 24
  issuer: go-ddd
```

入口接受 `-config`，默认值是 `config/config.yaml`。`config.Load` 实际调用的是 `viper.ReadConfig(os.Stdin)`，标准输入里的 YAML 会被解析；`-config` 指向的文件不会被读取。空的标准输入也会被当成一份空配置，随后补上这些缺省：端口 `8080`、模式 `debug`、空闲连接 `10`、最大连接 `100`、JWT 有效期 `24` 小时。

`Config.Database` 的 mapstructure 标签是 `databse`，对不上 YAML 里的 `database`。因此即使把上面的文件从标准输入送进去，数据库主机和端口仍是零值，DSN 会去拨号 `:0`。代码里还调用了 `viper.AutomaticEnv()`（配置键的 `.` 会换成 `_`），没有逐字段 `BindEnv`，环境变量能否覆盖配置未验证。

GORM 自动迁移只写在 `app.mode == "production"` 的分支里，迁移对象是 `UserModel`。Gin 的 `SetMode` 只接受 `debug`、`release`、`test`。把 `mode` 写成 `production` 会在设置 Gin 模式时 panic；写成 Gin 接受的值又不会进入这段自动迁移。`debug` / `test` 启动时不会建表。

## 环境要求

- Go 1.25.5（与 `go.mod` 一致；更低版本未验证）
- 可访问的 MySQL，库名与配置一致。当前示例是本机 `127.0.0.1:3306`，库 `go_ddd`，账号 `root` / `root`
- 本仓库的运行验证环境里没有 MySQL 客户端或服务

## 安装与运行

在仓库根目录：

```bash
go build -o /tmp/go-ddd ./cmd/api
```

本次在当前环境用 Go 1.25.5 执行上述命令，编译成功。

直接运行（标准输入为空，也没有 MySQL）：

```bash
./go-ddd -config config/config.yaml
```

本次实际输出：

```text
failed to init database: dial tcp :0: connect: connection refused
```

把 `config/config.yaml` 从标准输入送入时，仍然是同一条拨号错误（数据库配置没有绑定到结构体）。因此下面的接口调用在当前代码下还走不到处理函数，只列出与路由、DTO 一致的形状，供对照源码：

```bash
curl -s http://127.0.0.1:8080/heath

curl -s -X POST http://127.0.0.1:8080/api/v1/users/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"yiwen","email":"yiwen@example.com","password":"Admin123","nickname":"一文"}'

curl -s -X POST http://127.0.0.1:8080/api/v1/users/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"yiwen","password":"Admin123"}'
```

启动日志里的 “Api Base URL” 使用了已经带冒号的地址字符串，打印结果会是 `http://localhost::8080/api/v1` 这种形式。服务真正监听的是 `:<port>`。

### 数据库脚本

`scripts/sql/init.sql` 打算创建库 `go_ddd` 和表 `users`（自增 `id`、业务 `uuid`、用户名、邮箱、密码哈希、昵称、头像、状态、角色、创建/更新时间、软删除 `deleted_at`），并插入一条用户名为 `admin`、邮箱 `admin@example.com` 的记录。脚本注释写的密码是 `Admin123`。

这段脚本按文件内容不能直接执行：建表语句写成了 `CREAT TABLE`。注释中的 bcrypt 串也与 `Admin123` 对不上（已用 `bcrypt.CompareHashAndPassword` 核对过）。

## 和运行相关的现状

对照源码，这些行为会挡住一次完整的注册 / 登录：

1. 配置从标准输入读取，示例 YAML 的数据库段又因为结构体标签 `databse` 无法绑定，进程在监听端口之前就连接失败。
2. `scripts/sql/init.sql` 有语法错误；非 `production` 模式也不会 AutoMigrate。而 `production` 不是 Gin 接受的 mode。
3. `UserRepository.FindByUsername` 的方法体是 `panic("unimplemented")`。登录会调用它。
4. 个人资料、改密、按 ID 查询、`/me` 这组路由只挂了 `AdminMiddleware`，请求会以 403 结束。
5. 领域事件只保存在聚合对象里。封禁、提升管理员等命令类型没有对应的 HTTP 入口。

`UserModel` 的主键标签写成了 `primayKey`。若以后走到 AutoMigrate，GORM 对主键的识别以该标签为准。

## 参考

分层和命名对应的是常见 DDD 划分：实体与值对象、聚合根、仓储接口、领域服务、应用服务、接口层。读写模型拆在 `internal/application/command` 与 `internal/application/query`。实现范围以本仓库源码为准。
