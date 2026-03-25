package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SwordDragonLee/go-web/internal/api/handler"
	"github.com/SwordDragonLee/go-web/internal/api/router"
	"github.com/SwordDragonLee/go-web/internal/config"
	"github.com/SwordDragonLee/go-web/internal/domain/model"
	"github.com/SwordDragonLee/go-web/internal/repository"
	"github.com/SwordDragonLee/go-web/internal/service"
	"github.com/SwordDragonLee/go-web/internal/utils"
	"github.com/SwordDragonLee/go-web/pkg/db"
	"github.com/SwordDragonLee/go-web/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.InitLogger(cfg.Server.Mode); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	database, err := db.InitDB(&cfg.Database)
	if err != nil {
		logger.GetLogger().Fatal("初始化数据库失败", zap.Error(err))
		os.Exit(1)
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(
		&model.User{},
		&model.LoginLog{},
		&model.VerificationCode{},
		&model.Role{},
		&model.Permission{},
		&model.UserRoleLink{},
		&model.RolePermission{},
	); err != nil {
		logger.GetLogger().Fatal("数据库迁移失败", zap.Error(err))
		os.Exit(1)
	}

	// 配置JWT
	utils.SetJWTSecret(cfg.JWT.Secret)
	utils.SetTokenExpireDuration(cfg.JWT.ExpireTime)

	// 初始化仓储
	userRepo := repository.NewUserRepository(database)
	loginLogRepo := repository.NewLoginLogRepository(database)
	verificationCodeRepo := repository.NewVerificationCodeRepository(database)
	roleRepo := repository.NewRoleRepository(database)
	permissionRepo := repository.NewPermissionRepository(database)

	// 初始化服务
	userService := service.NewUserService(userRepo, loginLogRepo, verificationCodeRepo)
	roleService := service.NewRoleService(roleRepo)
	permissionService := service.NewPermissionService(permissionRepo)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)

	// 设置路由
	r := router.SetupRouter(userHandler, roleHandler, permissionHandler)

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器（在goroutine中）
	go func() {
		logger.GetLogger().Info("服务器启动", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.GetLogger().Info("正在关闭服务器...")

	// 优雅关闭，等待5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.GetLogger().Fatal("服务器强制关闭", zap.Error(err))
	}

	logger.GetLogger().Info("服务器已关闭")
}
