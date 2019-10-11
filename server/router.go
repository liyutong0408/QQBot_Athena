package server

import (
	"Athena/api"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由
func NewRouter() *gin.Engine {
	r := gin.Default()

	v := r.Group("/api")
	{
		v.POST("ReceiveMahuaOutput", api.Mahua)
	}

	return r
}
