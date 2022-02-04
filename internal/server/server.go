package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"newsplatform/internal/controllers"
)

func SetupServer() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": "pong",
		})
	})

	// Group API endpoints
	//r.Group("/api/v2")
	//{
	r.GET("/news", controllers.FindNews)
	r.POST("/news", controllers.CreateNews)
	r.GET("/news/:keyword", controllers.FindNew)
	r.PATCH("/news/:id", controllers.UpdateNew)
	r.DELETE("/news/:id", controllers.DeleteNews)
	r.DELETE("/newsbyday/:day", controllers.DeleteNewsByDay)
	//}
	return r
}
