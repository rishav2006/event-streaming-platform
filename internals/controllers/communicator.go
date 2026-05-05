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
	mu                  sync.Mutex
	LastFileNumOrder    int
	LastFileNumPayment  int
	Checker             bool
	OffsetSliceOrders   []int
	OffsetSlicePayments []int
	ExsOffset           int
}

// func offsetCount(c *gin.Context, path string, o *[]int, k int) {
// 	OffsetCalc(c, path, o, k)
// }

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

	// var offsetSliceOrders = []int{0, 0, 0}
	// var offsetSlicePayments = []int{0, 0}
	// var exsOffset int = 0

	// Check for checker, if it's false then the offset has not been calculated yet, so calculate it

	// path for orders requests
	var orders = []string{
		"internals/files/folders/orders/o0.log",
		"internals/files/folders/orders/o1.log",
		"internals/files/folders/orders/o2.log",
	}

	// path for payments requests
	var payments = []string{
		"internals/files/folders/payments/p0.log",
		"internals/files/folders/payments/p1.log",
	}

	// path for default requests
	var defaultPath string = "internals/files/folders/default/default.log"

	if d.Checker == false {
		// calculate offsets for all the files

		// For orders
		OffsetCalc(c, orders[0], &d.OffsetSliceOrders, 0)
		OffsetCalc(c, orders[1], &d.OffsetSliceOrders, 1)
		OffsetCalc(c, orders[2], &d.OffsetSliceOrders, 2)

		// For payments
		OffsetCalc(c, payments[0], &d.OffsetSlicePayments, 0)
		OffsetCalc(c, payments[1], &d.OffsetSlicePayments, 1)

		// For default
		OffsetCalcInt(c, defaultPath, &d.ExsOffset)

		d.Checker = true // set it to true
	}

	if stringOrders == true { // Order Case
		path := "internals/files/folders/orders"
		err := os.MkdirAll(path, 0755)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "error creating/opening the folder",
			})
			return
		}

		if d.LastFileNumOrder == 2 {
			file, _ := CheckFolder(orders[0], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSliceOrders[0], stringData)
			fmt.Println(line, 0)
			file.WriteString(line)
			d.OffsetSliceOrders[0]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to orders catalog",
			})

			file.Close()

		} else if d.LastFileNumOrder == 0 {
			file, _ := CheckFolder(orders[1], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSliceOrders[1], stringData)
			file.WriteString(line)
			fmt.Println(line, 1)
			d.OffsetSliceOrders[1]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to orders catalog",
			})

			file.Close()
		} else if d.LastFileNumOrder == 1 {
			file, _ := CheckFolder(orders[2], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSliceOrders[2], stringData)
			file.WriteString(line)
			fmt.Println(line, 2)
			d.OffsetSliceOrders[2]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to orders catalog",
			})

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

		if d.LastFileNumPayment == 0 {
			file, _ := CheckFolder(payments[1], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSlicePayments[1], stringData)
			file.WriteString(line)
			d.OffsetSlicePayments[1]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to payments catalog",
			})

			file.Close()
		} else if d.LastFileNumPayment == 1 {
			file, _ := CheckFolder(payments[0], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSlicePayments[0], stringData)
			file.WriteString(line)
			d.OffsetSlicePayments[0]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to payments catalog",
			})

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
		// f, err := os.Open("internals/files/folders/default/default.log")
		// if err != nil {
		// 	c.JSON(400, gin.H{
		// 		"error": "there was some error oprning the file for offset counting",
		// 	})
		// 	return
		// }

		// scanner := bufio.NewScanner(f)

		// for scanner.Scan() {
		// 	exsOffset++
		// }

		// defer f.Close()

		line := fmt.Sprintf("%d | %s\n", d.ExsOffset, stringData)
		d.ExsOffset++
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
