package router

import (
	"github.com/gin-gonic/gin"
	"slot-crawler/internal/server/router/handler"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	apiRouter := r.Group("/api")
	{
		apiRouter.GET("/users/:slot", handler.GetUserList)
		apiRouter.GET("/data", handler.GetSpinData)
		apiRouter.GET("/crawl", handler.StartCrawling)
	}
	return r
}
