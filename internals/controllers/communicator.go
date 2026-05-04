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
	mu                 sync.Mutex
	LastFileNumOrder   int
	LastFileNumPayment int
}

func (d *Demo) Producer(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error reading the contents of the file",
		})
		return
	}
	stringData := string(jsonData) // Get the data from the request.body

	_, stringOrders := c.GetQuery("order")     // Get the optional parameters 1
	_, stringPayments := c.GetQuery("payment") // Get the optional parameters 2

	d.mu.Lock()
	defer d.mu.Unlock()

	if stringOrders == true { // Order Case
		path := "internals/files/folders/orders"
		err := os.MkdirAll(path, 0755)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "error creating/opening the folder",
			})
			return
		}
		var offsetSliceOrders = []int{0, 0, 0}

		if d.LastFileNumOrder == 2 {
			file, err := os.OpenFile("internals/files/folders/orders/o0.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				c.JSON(400, gin.H{
					"error": "error creating/opening the file",
				})
			}
			f, err := os.Open("internals/files/folders/orders/o0.log")
			if err != nil {
				c.JSON(400, gin.H{
					"error": "There was some error opening the file",
				})
				return
			}
			scanner := bufio.NewScanner(f)
			defer f.Close()

			for scanner.Scan() {
				offsetSliceOrders[0]++
			}

			line := fmt.Sprintf("%d | %s\n", offsetSliceOrders[0], stringData)
			file.WriteString(line)
			offsetSliceOrders[0]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to orders catalog",
			})
			file.Close()

		} else if d.LastFileNumOrder == 0 {
			file, err := os.OpenFile("internals/files/folders/orders/o1.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				c.JSON(400, gin.H{
					"error": "error creating/opening the file",
				})
				return
			}
			f, err := os.Open("internals/files/folders/orders/o1.log")
			if err != nil {
				c.JSON(400, gin.H{
					"error": "There was some error opening the file",
				})
				return
			}
			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				offsetSliceOrders[1]++
			}

			line := fmt.Sprintf("%d | %s\n", offsetSliceOrders[1], stringData)
			file.WriteString(line)
			offsetSliceOrders[1]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to orders catalog",
			})
			file.Close()
		} else if d.LastFileNumOrder == 1 {
			file, err := os.OpenFile("internals/files/folders/orders/o2.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				c.JSON(400, gin.H{
					"error": "error creating/opening the file",
				})
				return
			}

			f, err := os.Open("internals/files/folders/orders/o2.log")
			if err != nil {
				c.JSON(400, gin.H{
					"error": "There was some error opening the file",
				})
				return
			}
			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				offsetSliceOrders[2]++
			}

			line := fmt.Sprintf("%d | %s\n", offsetSliceOrders[2], stringData)
			file.WriteString(line)
			offsetSliceOrders[2]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to orders catalog",
			})
			f.Close()
			file.Close()
		}
		d.LastFileNumOrder = (d.LastFileNumOrder + 1) % 3

	} else if stringPayments == true { // Payment case
		path := "internals/files/folders/payments"
		err := os.MkdirAll(path, 0755)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "error creating/opening the folder",
			})
			return
		}
		var offsetSlicePayments = []int{0, 0}

		if d.LastFileNumPayment == 0 {
			file, err := os.OpenFile("internals/files/folders/payments/p1.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				c.JSON(400, gin.H{
					"error": "error creating/opening the file",
				})
				return
			}
			f, err := os.Open("internals/files/folders/payments/p1.log")
			if err != nil {
				c.JSON(400, gin.H{
					"error": "There was some error opening the file",
				})
				return
			}
			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				offsetSlicePayments[1]++
			}

			line := fmt.Sprintf("%d | %s\n", offsetSlicePayments[1], stringData)
			file.WriteString(line)
			offsetSlicePayments[1]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to payments catalog",
			})
			f.Close()
			file.Close()
		} else if d.LastFileNumPayment == 1 {
			file, err := os.OpenFile("internals/files/folders/payments/p0.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				c.JSON(400, gin.H{
					"error": "error creating/opening the file",
				})
				return
			}
			f, err := os.Open("internals/files/folders/payments/p0.log")
			if err != nil {
				c.JSON(400, gin.H{
					"error": "There was some error opening the file",
				})
				return
			}
			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				offsetSlicePayments[0]++
			}

			line := fmt.Sprintf("%d | %s\n", offsetSlicePayments[0], stringData)
			file.WriteString(line)
			offsetSlicePayments[0]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to payments catalog",
			})
			f.Close()
			file.Close()
		}
		d.LastFileNumPayment = (d.LastFileNumPayment + 1) % 2

	} else if stringOrders == false && stringPayments == false { // Default case
		path := "internals/files/folders/default"
		err := os.MkdirAll(path, 0755)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "error creating/opening the folder",
			})
			return
		}

		// Store it in file
		file, err := os.OpenFile("internals/files/folders/default/default.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "error opening the file",
			})
			return
		}

		defer file.Close()

		// Find the offset
		f, err := os.Open("internals/files/folders/default/default.log")
		if err != nil {
			c.JSON(400, gin.H{
				"error": "there was some error oprning the file for offset counting",
			})
			return
		}
		scanner := bufio.NewScanner(f)
		var exsOffset int = 0
		for scanner.Scan() {
			exsOffset++
		}

		defer f.Close()

		line := fmt.Sprintf("%d | %s\n", exsOffset, stringData)
		exsOffset++
		file.WriteString(line)

		c.JSON(201, gin.H{
			"message": "message sent successfully",
		})
	}

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
