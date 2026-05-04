package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rishav2006/event-streaming/internals/controllers"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	var d = controllers.Demo{LastFileNumOrder: 2, LastFileNumPayment: 1}

	r.POST("/produce", d.Producer)
	r.GET("/consume", d.Consumer)

	return r;
}