package router

import (
	"collectify/internal/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	initCategoryRouter(router)

	return router
}

func initCategoryRouter(router *gin.Engine) {
	category := router.Group("/category")
	{
		category.POST("", handler.CreateCategory)
	}
}
