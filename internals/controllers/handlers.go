package controllers

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func OffsetCalc2(c *gin.Context, path string, s *[]int, k int) {
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

func OffsetCalc(c *gin.Context, path string, t *int, s *[]int, k int) {
	f, err := os.Open(path)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "there was some error opening the file for offset calculation",
		})
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		*t++
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

func ConsumerReadFiles(c *gin.Context, path string, offset int, groupString string, mpp map[string][]int) ([]Answer, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error reading/listing the files from the specified directory",
		})
		return []Answer{}, err
	}

	var AnswerArray = make([]Answer, 0, 10)

	var i int = 1
	var e int = 0
	var arr = mpp[groupString]
	var line string

	for k := 0; k < len(arr); k++ {
		for arr[k] > 0 {
			var count int = 0
			entry := entries[e]
			e++
			fullPath := filepath.Join(path, entry.Name())
			file, err := os.Open(fullPath)
			if err != nil {
				c.JSON(400, gin.H{
					"error": "could not open the file",
				})
				return []Answer{}, err
			}
			scanner := bufio.NewScanner(file)

			var conNum string = "Consumer " + strconv.Itoa(i)
			i++

			for scanner.Scan() {
				if count < offset {
					count++
				} else {
					// var offsetNum string
					line = scanner.Text()
					re := regexp.MustCompile(`^[\d]+`)
					match := re.FindString(line)
					_, actualMessage, found := strings.Cut(line, "| ")
					if found {
						AnswerArray = append(AnswerArray, Answer{Group: groupString, Consumer: conNum, Message: actualMessage, Partition: entry.Name(), Offset: match})
					}
				}
			}
			arr[k]--
		}
	}

	return AnswerArray, nil
}

func ConsumerReadFilesNoGroup(c *gin.Context, path string, offset int, topic string) ([]AnswerNoGroup, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error reading/listing the files from the specified directory",
		})
		return []AnswerNoGroup{}, err
	}

	var AnswerArray = make([]AnswerNoGroup, 0, 10)

	for _, entry := range entries {
		var count int = 0

		fullPath := filepath.Join(path, entry.Name())
		file, err := os.Open(fullPath)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "could not open the file",
			})
			return []AnswerNoGroup{}, err
		}
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			if count < offset {
				count++
			} else {
				line := scanner.Text()
				AnswerArray = append(AnswerArray, AnswerNoGroup{Message: line, Topic: topic, Partition: entry.Name()})
			}
		}

	}
	return AnswerArray, nil
}

func PartitionNumFinder(path string, topic string) string {
	if topic == "orders" {
		var s1 = strings.ReplaceAll(path, "internals/files/folders/orders/", "")
		var s2 = strings.ReplaceAll(s1, ".log", "")
		return s2
	} else if topic == "payments" {
		var s1 = strings.ReplaceAll(path, "internals/files/folders/payments/", "")
		var s2 = strings.ReplaceAll(s1, ".log", "")
		return s2
	} else {
		var s1 = strings.ReplaceAll(path, "internals/files/folders/default/", "")
		var s2 = strings.ReplaceAll(s1, ".log", "")
		return s2
	}
}
