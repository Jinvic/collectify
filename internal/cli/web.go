package cli

import (
	"collectify/internal/config"
	"collectify/internal/router"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func DoWeb() {
	cfg := config.GetConfig()

	// 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化路由
	r := router.InitRouter()

	// 配置 CORS for development
	// In production, you might want to restrict this more or remove it
	// if the frontend is served from the same origin.
	// This should ideally be configurable via config file.
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"} // React dev server
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Serve frontend static files
	// Check if the frontend build directory exists
	frontendBuildPath := "./web/build"
	if _, err := os.Stat(frontendBuildPath); err == nil {
		// Serve static files (CSS, JS, images)
		r.Static("/static", filepath.Join(frontendBuildPath, "static"))

		// Serve favicon
		r.StaticFile("/favicon.ico", filepath.Join(frontendBuildPath, "favicon.ico"))

		// Catch-all handler for SPA (Single Page Application)
		// This will serve index.html for any route that is not an API call or a static asset.
		// This allows React Router to handle routing on the client side.
		r.NoRoute(func(c *gin.Context) {
			// If the request is for an API or a known static asset, return 404
			if strings.HasPrefix(c.Request.URL.Path, "/api/") ||
				strings.HasPrefix(c.Request.URL.Path, "/static/") ||
				c.Request.URL.Path == "/favicon.ico" {
				c.AbortWithStatus(404)
				return
			}
			// For all other routes, serve the React index.html file.
			c.File(filepath.Join(frontendBuildPath, "index.html"))
		})
	} else {
		// If frontend is not built, log a message
		log.Printf("Frontend build directory '%s' not found. Serving API only.\n", frontendBuildPath)
		// You might want to serve a simple API-only page or just let the API routes handle everything.
		// For now, we'll just log and let the API routes function normally.
	}

	// 创建 HTTP Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 启动服务器（goroutine）
	go func() {
		log.Printf("🚀 Server is running at http://localhost:%d\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server startup failed: %v\n", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("_shutdown signal received...")

	// 创建超时上下文，用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 优雅关闭服务器（等待正在处理的请求完成）
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v\n", err)
	}

	log.Println("✅ Server exited gracefully")
}