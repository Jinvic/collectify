package router

import (
	"collectify/internal/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 将所有 API 路由分组到 /api 路径下
	api := router.Group("/api")
	{
		initCategoryRouter(api)
		initFieldRouter(api)
		initItemRouter(api)
	}

	// Serve frontend static files
	// This part is handled in internal/cli/web.go

	return router
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