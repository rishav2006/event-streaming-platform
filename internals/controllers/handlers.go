package controllers

import (
	"bufio"
	"os"
	"path/filepath"

	// "text/scanner"

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

func ConsumerReadFiles(c *gin.Context, path string) (string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error reading/listing the files from the specified directory",
		})
		return "", err
	}

	var str string

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		file, err := os.Open(fullPath)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "could not open the file",
			})
			return "", err
		}
		scanner := bufio.NewScanner(file)
		// var i int = 0
		for scanner.Scan() {
			line := scanner.Text()
			str = str + line
			str += "\n"
		}
	}
	return str, nil
}
