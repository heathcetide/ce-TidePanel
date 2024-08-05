package router

import (
	"github.com/gin-gonic/gin"
)

func CollectionRoute(r *gin.Engine) *gin.Engine {
	//r.Use(gin.LoggerWithConfig(logger.LoggerToFile()))
	//r.Use(logger.Recover)
	//r.POST("/judge", handler.HandlerJudgeQuestion)
	//r.POST("/addQuestion", handler.HandlerAddQuestion)
	r.GET("/docs", GetHtml)
	r.GET("/api/data/:tableName", GetTableData)
	r.GET("/api/table/:tableName", GetTableInfo)
	return r
}
