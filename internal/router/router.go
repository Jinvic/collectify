package router

import (
	"collectify/internal/handler"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

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
	fsys, err := fs.Sub(frontendEmbedFS, "web/build")
	if err != nil {
		return nil
	}
	return http.FS(fsys)
}
func getFrontendStaticFS() http.FileSystem {
	fsys, err := fs.Sub(frontendEmbedFS, "web/build/static")
	if err != nil {
		return nil
	}
	return http.FS(fsys)
}


func InitRouter() *gin.Engine {
	r := gin.Default()

	// 配置 CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// 获取前端静态文件
	frontendFS := getFrontendFS()
	frontendStaticFS := getFrontendStaticFS()
	if frontendFS != nil && frontendStaticFS != nil {
		log.Println("✅ 已加载嵌入的前端静态文件")
	} else {
		log.Println("⚠️ 未找到嵌入的前端文件，仅提供 API 服务")
	}

	if frontendFS != nil {
		// 挂载静态目录
		r.StaticFS("/static", frontendStaticFS)

		// favicon
		r.GET("/favicon.ico", func(c *gin.Context) {
			c.FileFromFS("favicon.ico", frontendFS)
		})

		// 根路径返回 index.html
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

	// 将所有 API 路由分组到 /api 路径下
	api := r.Group("/api")
	{
		initCategoryRouter(api)
		initFieldRouter(api)
		initItemRouter(api)
	}

	return r
}

func initCategoryRouter(router *gin.RouterGroup) {
	category := router.Group("/category")
	{
		category.POST("", handler.CreateCategory)
		category.GET("/:id", handler.GetCategory)
		category.PATCH("/:id", handler.RenameCategory)
		category.DELETE("/:id", handler.DeleteCategory)
		category.GET("/list", handler.ListCategory)
		category.POST("/:id/restore", handler.RestoreCategory)
	}
}

func initFieldRouter(router *gin.RouterGroup) {
	field := router.Group("/field")
	{
		field.POST("", handler.CreateField)
		field.DELETE("/:id", handler.DeleteField)
		field.POST("/:id/restore", handler.RestoreField)
	}
}

func initItemRouter(router *gin.RouterGroup) {
	item := router.Group("/item")
	{
		item.POST("", handler.CreateItem)
		item.DELETE("/:id", handler.DeleteItem)
		item.PUT("/:id", handler.UpdateItem)
		item.POST("/:id/restore", handler.RestoreItem)
		item.GET("/list", handler.ListItems)
		item.POST("/search", handler.SearchItems)
		item.GET("/:id", handler.GetItem)
	}
}
