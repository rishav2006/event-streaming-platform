package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rishav2006/event-streaming/internals/controllers"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/produce", controllers.Producer)
	r.GET("/consume", controllers.Consumer)

	return r;
}