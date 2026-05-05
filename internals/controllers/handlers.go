package controllers

import (
	"bufio"
	"os"

	"github.com/gin-gonic/gin"
)

func OffsetCalc(c *gin.Context, path string, s *[]int, k int) {
	f, err := os.Open(path)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "there was some error opening the file for offset calculation",
		})
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		(*s)[k]++
	}
	f.Close()
}

func OffsetCalcInt(c *gin.Context, path string, excs *int) {
	f, err := os.Open(path)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "there was some error opening the file for offset calculation",
		})
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		*excs++
	}
	f.Close()
}

func CheckFolder(path string, c *gin.Context) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error opening/creating the folder",
		})
		return nil, err
	}
	return file, nil
}
