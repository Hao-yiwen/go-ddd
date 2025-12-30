package main

import (
	"flag"
	"fmt"
	"log"
	"yiwen/go-ddd/internal/application/service"
	"yiwen/go-ddd/internal/infrastructure/config"
	"yiwen/go-ddd/internal/infrastructure/persistence/model"
	"yiwen/go-ddd/internal/interfaces/api/handler"
	"yiwen/go-ddd/internal/interfaces/api/middleware"
	"yiwen/go-ddd/internal/interfaces/api/router"

	mysqlrepo "yiwen/go-ddd/internal/infrastructure/persistence/mysql"

	domainservice "yiwen/go-ddd/internal/domain/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	gin.SetMode(cfg.App.Mode)

	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	userRepo := mysqlrepo.NewUserRepository(db)

	userDomainService := domainservice.NewUserDomainService(userRepo)

	userApplicationService := service.NewUserApplicationService(userRepo, *userDomainService)

	jwtAuth := middleware.NewJWTAuth(cfg.JWT.Secret, cfg.JWT.ExpireHour, cfg.JWT.Issuer)

	userHandler := handler.NewUserHandler(userApplicationService, jwtAuth)

	r := router.NewRouter(userHandler, jwtAuth)

	engine := r.Setup()

	addr := fmt.Sprintf(":%d", cfg.App.Port)

	log.Printf("server is running on port %d", cfg.App.Port)
	log.Printf("mode: %s", cfg.App.Mode)
	log.Printf("Api Base URL: http://localhost:%s/api/v1", addr)

	if err := engine.Run(addr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	var logLevel logger.LogLevel
	switch cfg.App.Mode {
	case "debug":
		logLevel = logger.Info
	case "test":
		logLevel = logger.Warn
	case "production":
		logLevel = logger.Error
	default:
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)

	// 自动迁移（AutoMigrate）是 GORM 提供的一种功能，用于自动根据模型（如 UserModel）的结构变更数据库中的表结构。
	// 它会自动创建、修改（字段类型默认兼容）、删除表结构，以确保数据库结构和 Go 代码中的模型结构保持一致。
	// 但是，自动迁移不会删除已有字段的内容，也不会更改已有字段的类型，它主要用于保持字段的增加或表的创建同步。
	// 这里示例只在生产环境（production）下才自动迁移，以避免开发或测试时误操作数据库结构。
	if cfg.App.Mode == "production" {
		if err := db.AutoMigrate(&model.UserModel{}); err != nil {
			return nil, err
		}
	}

	return db, nil

}
