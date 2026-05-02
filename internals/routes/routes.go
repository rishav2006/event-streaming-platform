package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rishav2006/event-streaming/internals/controllers"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/send", controllers.Sender)
	r.GET("/consume", controllers.ReceiverOffset)

	return r;
}