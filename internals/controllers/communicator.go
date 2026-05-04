package controllers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Demo struct {
	mu sync.Mutex
}

func (d *Demo) Producer(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error reading the contents of the file",
		})
		return
	}
	stringData := string(jsonData)

	d.mu.Lock()
	defer d.mu.Unlock()

	// Store it in file
	file, err := os.OpenFile("internals/files/test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error opening the file",
		})
		return
	}

	defer file.Close()

	// Find the offset 
	f, err := os.Open("internals/files/test.log")
	if err != nil {
		c.JSON(400, gin.H{
			"error" : "there was some error oprning the file for offset counting",
		})
	}
	scanner := bufio.NewScanner(f)
	var exsOffset int = 0;
	for scanner.Scan() {
		exsOffset++;
	}

	defer f.Close()

	line := fmt.Sprintf("%d | %s\n", exsOffset, stringData)
	exsOffset++
	file.WriteString(line)

	c.JSON(201, gin.H{
		"message": "message sent successfully",
	})

}

func (d *Demo) Consumer(c *gin.Context) {
	file, err := os.Open("internals/files/test.log")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "There was some error opening the file",
		})
		return
	}
	scanner := bufio.NewScanner(file)
	defer file.Close()

	d.mu.Lock()
	defer d.mu.Unlock()

	offsetString := c.Query("offset")
	if offsetString == "" {
		// If offset is not present, scan through the entire file
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
		return
	}
	var cnt int = 0
	// If offset is present, start scanning from the offset position of the file
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