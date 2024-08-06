package router

import (
	"github.com/gin-gonic/gin"
)

func CollectionRoute(r *gin.Engine) *gin.Engine {
	r.GET("/docs", GetHtml)
	r.GET("/api/data/:tableName", GetTableData)
	r.GET("/api/table/:tableName", GetTableInfo)
	return r
}
