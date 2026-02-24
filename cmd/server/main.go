// GoCargo 库存管理系统入口
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-cargo/internal/config"
	"go-cargo/internal/database"
	"go-cargo/internal/handler"
	"go-cargo/internal/repository"
	"go-cargo/internal/router"
	"go-cargo/internal/service"
	"go-cargo/web"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 设置 Gin 模式
	if cfg.AppMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	db := database.Init(cfg)

	// 初始化各层
	repo := repository.New(db)
	svc := service.New(repo, cfg)
	h := handler.New(svc)

	// 设置路由
	r := router.Setup(h, web.StaticFS)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppPort),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器 (非阻塞)
	go func() {
		log.Printf("🚀 GoCargo 库存管理系统已启动")
		log.Printf("📍 访问地址: http://localhost:%s", cfg.AppPort)
		log.Printf("📍 API 地址: http://localhost:%s/api/v1", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}
	log.Println("服务器已安全关闭")
}
