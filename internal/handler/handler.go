package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SuccessCode = 0
	FailCode    = 1
)

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": SuccessCode,
		"msg":  "success",
	})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": SuccessCode,
		"msg":  "success",
		"data": data,
	})
}

func Fail(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code": FailCode,
		"msg":  err.Error(),
	})
}

func FailWithData(c *gin.Context, err error, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": FailCode,
		"msg":  err.Error(),
		"data": data,
	})
}
