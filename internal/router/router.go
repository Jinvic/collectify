package router

import (
	"collectify/internal/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	initCategoryRouter(router)
	initFieldRouter(router)

	initDeletedRouter(router)

	return router
}

func initCategoryRouter(router *gin.Engine) {
	category := router.Group("/category")
	{
		category.POST("", handler.CreateCategory)
		category.GET("/:id", handler.GetCategory)
		category.PATCH("/:id", handler.RenameCategory)
		category.DELETE("/:id", handler.DeleteCategory)
		category.POST("/search", handler.SearchCategory)
	}
}

func initDeletedRouter(router *gin.Engine) {
	deleted := router.Group("/deleted")
	{
		deleted.POST("", handler.Restore)
		deleted.DELETE("", handler.HardDelete)
	}
}

func initFieldRouter(router *gin.Engine) {
	field := router.Group("/field")
	{
		field.POST("", handler.CreateField)
		field.DELETE("/:id", handler.DeleteField)
	}
}
