package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rishav2006/event-streaming/internals/models"
)

func Sender(c *gin.Context) {
	var newSender models.EventModel
	newSender.Offset = models.Count
	models.Count = models.Count + 1
	if err := c.BindJSON(&newSender); err != nil {
		c.JSON(400, gin.H{
			"error": "error sending the message",
		})
		return
	}
	models.Events = append(models.Events, newSender)
	c.JSON(201, gin.H{
		"message": "message sent successfully",
	})
}

func Receiver(c *gin.Context) {
	c.JSON(200, models.Events)
}

func ReceiverOffset(c *gin.Context) {
	var offsetStr = c.Query("offset")
	if offsetStr == "" {
		c.JSON(200, models.Events)
		return
	}
	var offset, err = strconv.Atoi(offsetStr);
	if err != nil {
		c.JSON(400, gin.H{
			"error" : "failed to convert from string to integer",
		})
	}

	for i, event := range models.Events {
		if i >= offset {
			c.JSON(200, event);
		}
	}
}
