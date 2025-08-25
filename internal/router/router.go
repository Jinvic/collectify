package router

import (
	"collectify/internal/handler"
	"collectify/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 配置 CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// 注册 API 路由
	api := r.Group("/api")
	{
		initCategoryRouter(api)
		initCollectionRouter(api)
		initFieldRouter(api)
		initItemRouter(api)
		initTagRouter(api)
		initUserRouter(api)
	}

	// 初始化前端路由
	initFrontendRouter(r)

	return r
}

func initCategoryRouter(router *gin.RouterGroup) {
	category := router.Group("/category")
	{
		category.GET("/:id", handler.GetCategory)
		category.GET("/list", handler.ListCategory)

		category.Use(middleware.AuthCheck)
		category.POST("", handler.CreateCategory)
		category.PATCH("/:id", handler.RenameCategory)
		category.DELETE("/:id", handler.DeleteCategory)
		category.POST("/:id/restore", handler.RestoreCategory)
	}
}

func initCollectionRouter(router *gin.RouterGroup) {
	collection := router.Group("/collection")
	{
		collection.GET("/:id", handler.GetCollection)
		collection.GET("/list", handler.ListCollection)

		collection.Use(middleware.AuthCheck)
		collection.POST("", handler.CreateCollection)
		collection.PATCH("/:id", handler.UpdateCollection)
		collection.DELETE("/:id", handler.DeleteCollection)
		collection.POST("/:id/restore", handler.RestoreCollection)
	}
}

func initFieldRouter(router *gin.RouterGroup) {
	field := router.Group("/field")
	{
		field.Use(middleware.AuthCheck)
		field.POST("", handler.CreateField)
		field.DELETE("/:id", handler.DeleteField)
		field.POST("/:id/restore", handler.RestoreField)
	}
}

func initItemRouter(router *gin.RouterGroup) {
	item := router.Group("/item")
	{
		item.GET("/list", handler.ListItems)
		item.POST("/search", handler.SearchItems)
		item.GET("/:id", handler.GetItem)

		item.Use(middleware.AuthCheck)
		item.POST("", handler.CreateItem)
		item.DELETE("/:id", handler.DeleteItem)
		item.PUT("/:id", handler.UpdateItem)
		item.POST("/:id/restore", handler.RestoreItem)

		// 关联关系
		item.POST("/:id/tag/:tag_id", handler.AddTag)
		item.DELETE("/:id/tag/:tag_id", handler.RemoveTag)
		item.POST("/:id/collection/:collection_id", handler.AddToCollection)
		item.DELETE("/:id/collection/:collection_id", handler.RemoveFromCollection)
	}
}

func initTagRouter(router *gin.RouterGroup) {
	tag := router.Group("/tag")
	{
		tag.GET("/:id", handler.GetTag)
		tag.GET("/list", handler.ListTag)

		tag.Use(middleware.AuthCheck)
		tag.POST("", handler.CreateTag)
		tag.PATCH("/:id", handler.RenameTag)
		tag.DELETE("/:id", handler.DeleteTag)
		tag.POST("/:id/restore", handler.RestoreTag)
	}
}

func initUserRouter(router *gin.RouterGroup) {
	user := router.Group("/user")
	{
		user.POST("/login", handler.UserLogin)
		user.GET("/enabled", handler.IsAuthEnabled)

		user.Use(middleware.AuthCheck)
		user.POST("/update", handler.UserUpdate)
	}
}
