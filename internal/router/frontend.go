package router

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var frontendEmbedFS embed.FS

// 设置前端文件
func SetFrontendFS(fs embed.FS) {
	frontendEmbedFS = fs
}

// 将 embed.FS 转换为 http.FileSystem
func getFrontendFS() http.FileSystem {
	fsys, err := fs.Sub(frontendEmbedFS, "web/build")
	if err != nil {
		return nil
	}
	return http.FS(fsys)
}

func initFrontendRouter(r *gin.Engine) {
	// 获取前端文件系统
	frontendFS := getFrontendFS()
	if frontendFS != nil {
		log.Println("✅ 已加载嵌入的前端静态文件")

		// 手动处理文件
		fileHandler := createFileHandler(frontendFS)

		// 处理静态文件
		r.GET("/static/*filepath", func(c *gin.Context) {
			filepath := strings.TrimPrefix(c.Param("filepath"), "/")
			fullPath := "static/" + filepath
			fileHandler(c, fullPath)
		})

		// 处理根路径
		r.GET("/", func(c *gin.Context) {
			fileHandler(c, "index.html")
		})

		// 处理根路径文件
		r.GET("/index.html", func(c *gin.Context) {
			fileHandler(c, "index.html")
		})
		r.GET("/favicon.ico", func(c *gin.Context) {
			fileHandler(c, "favicon.ico")
		})
		r.GET("/robots.txt", func(c *gin.Context) {
			fileHandler(c, "robots.txt")
		})
		r.GET("/manifest.json", func(c *gin.Context) {
			fileHandler(c, "manifest.json")
		})
		r.GET("/asset-manifest.json", func(c *gin.Context) {
			fileHandler(c, "asset-manifest.json")
		})
		r.GET("/logo192.png", func(c *gin.Context) {
			fileHandler(c, "logo192.png")
		})

		// SPA 路由兜底 - 处理所有前端路由
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			// 排除 API 和静态文件
			if strings.HasPrefix(path, "/api/") ||
				strings.HasPrefix(path, "/static/") {
				c.AbortWithStatus(404)
				return
			}

			// 所有前端路由返回 index.html
			fileHandler(c, "index.html")
			// c.FileFromFS("index.html", frontendFS) // Cause infinite redirect
		})
	} else {
		log.Println("⚠️ 未找到嵌入的前端文件，仅提供 API 服务")
	}

}

// 文件处理器
func createFileHandler(fs http.FileSystem) func(c *gin.Context, path string) {
	return func(c *gin.Context, path string) {
		// 打开文件
		file, err := fs.Open(path)
		if err != nil {
			log.Printf("文件未找到: %s, 错误: %v", path, err)
			c.AbortWithStatus(404)
			return
		}
		defer file.Close()

		// 获取文件信息
		stat, err := file.Stat()
		if err != nil {
			log.Printf("文件状态错误: %s, 错误: %v", path, err)
			c.AbortWithStatus(404)
			return
		}

		// 设置正确的 Content-Type
		ext := filepath.Ext(path)
		switch ext {
		case ".css":
			c.Header("Content-Type", "text/css")
		case ".js":
			c.Header("Content-Type", "application/javascript")
		case ".html":
			c.Header("Content-Type", "text/html")
		case ".ico":
			c.Header("Content-Type", "image/x-icon")
		case ".png":
			c.Header("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			c.Header("Content-Type", "image/jpeg")
		case ".svg":
			c.Header("Content-Type", "image/svg+xml")
		case ".json":
			c.Header("Content-Type", "application/json")
		default:
			// 让浏览器自动判断类型
		}

		// 设置缓存头 - 对静态资源有效
		if strings.HasPrefix(path, "static/") {
			c.Header("Cache-Control", "public, max-age=31536000")
		}

		// 使用 http.ServeContent 提供文件
		http.ServeContent(c.Writer, c.Request, path, stat.ModTime(), file)
	}
}
