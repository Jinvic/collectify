package cli

import (
	"collectify/internal/config"
	"collectify/internal/router"
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"embed"

	"github.com/gin-gonic/gin"
)

var frontendEmbedFS embed.FS

// 设置前端文件
func SetFrontendFS(fs embed.FS) {
	frontendEmbedFS = fs
}

// 将 embed.FS 转换为 http.FileSystem
func getFrontendFS() http.FileSystem {
	if fsys, err := fs.Sub(frontendEmbedFS, "web/build"); err == nil {
		return http.FS(fsys)
	}
	return nil
}

func DoWeb() {
	cfg := config.GetConfig()

	// 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 获取前端静态文件
	frontendFS := getFrontendFS()
	if frontendFS != nil {
		log.Println("✅ 已加载嵌入的前端静态文件")
	} else {
		log.Println("⚠️ 未找到嵌入的前端文件，仅提供 API 服务")
	}
	r := router.InitRouter(frontendFS)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Printf("🚀 服务已启动：http://localhost:%d\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ 服务启动失败：%v\n", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("收到关闭信号...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ 服务强制关闭：%v\n", err)
	}

	log.Println("✅ 服务已退出")
}
