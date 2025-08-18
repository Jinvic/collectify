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
	"strings"
	"syscall"
	"time"

	"embed"

	"github.com/gin-contrib/cors"
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

	r := router.InitRouter()

	// 配置 CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// 尝试使用嵌入的前端文件系统
	var frontendFS http.FileSystem
	var useFrontend bool

	if fs := getFrontendFS(); fs != nil {
		frontendFS = fs
		useFrontend = true
		log.Println("✅ 已加载嵌入的前端静态文件")
	} else {
		log.Println("⚠️ 未找到嵌入的前端文件，仅提供 API 服务")
	}

	if useFrontend {
		// 提供静态资源
		r.StaticFS("/static", frontendFS)

		// 提供 favicon
		r.GET("/favicon.ico", func(c *gin.Context) {
			c.FileFromFS("/favicon.ico", frontendFS)
		})

		// 显式处理根路径
		r.GET("/", func(c *gin.Context) {
			c.FileFromFS("index.html", frontendFS)
		})

		// SPA 路由兜底
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api/") ||
				strings.HasPrefix(path, "/static/") ||
				path == "/favicon.ico" {
				c.AbortWithStatus(404)
				return
			}
			c.FileFromFS("index.html", frontendFS)
		})
	}

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
