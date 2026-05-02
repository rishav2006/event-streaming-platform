package controllers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Producer(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error reading the contents of the file",
		})
		return
	}
	stringData := string(jsonData)

	// Store it in file
	file, err := os.OpenFile("internals/files/test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error opening the file",
		})
	}

	// Find the offset 
	f, err := os.Open("internals/files/test.log")
	scanner := bufio.NewScanner(f)
	var exsOffset int = 0;
	for scanner.Scan() {
		exsOffset++;
	}

	line := fmt.Sprintf("%d | %s\n", exsOffset, stringData)
	exsOffset++
	file.WriteString(line)
	c.JSON(201, gin.H{
		"message": "message sent successfully",
	})
}

func Consumer(c *gin.Context) {
	file, err := os.Open("internals/files/test.log")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "There was some error opening the file",
		})
	}
	scanner := bufio.NewScanner(file)
	defer file.Close()

	offsetString := c.Query("offset")
	if offsetString == "" {
		// Scan through the file and find out
		for scanner.Scan() {
			line := scanner.Text()
			c.JSON(200, gin.H{
				"message": line,
			})
		}
		return
	}
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "failed to convert from string to integer",
		})
	}
	var cnt int = 0
	// Scan through the file and find out
	for scanner.Scan() {
		if cnt >= offset {
			line := scanner.Text()
			c.JSON(200, gin.H{
				"message": line,
			})
		}
		cnt++
	}
}
