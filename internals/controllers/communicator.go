package controllers

import (
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
	CounterOrder        int
	CounterPayment      int
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

		if d.CounterOrder%3 == 0 {
			file, _ := CheckFolder(orders[0], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSliceOrders[0], stringData)
			fmt.Println(line, 0)
			file.WriteString(line)
			d.OffsetSliceOrders[0]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to orders catalog",
			})
			file.Close()

		} else if d.CounterOrder%3 == 1 {
			file, _ := CheckFolder(orders[1], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSliceOrders[1], stringData)
			file.WriteString(line)
			fmt.Println(line, 1)
			d.OffsetSliceOrders[1]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to orders catalog",
			})
			file.Close()

		} else if d.CounterOrder%3 == 2 {
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
		d.CounterOrder++

	} else if stringPayments == true { // Payment case
		path := "internals/files/folders/payments"
		err := os.MkdirAll(path, 0755)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "error creating/opening the folder",
			})
			return
		}

		if d.CounterPayment%2 == 1 {
			file, _ := CheckFolder(payments[1], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSlicePayments[1], stringData)
			file.WriteString(line)
			d.OffsetSlicePayments[1]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to payments catalog",
			})

			file.Close()
		} else if d.CounterPayment%2 == 0 {
			file, _ := CheckFolder(payments[0], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSlicePayments[0], stringData)
			file.WriteString(line)
			d.OffsetSlicePayments[0]++
			c.JSON(201, gin.H{
				"message": "message successfully sent to payments catalog",
			})

			file.Close()
		}
		d.CounterPayment++

	} else if stringOrders == false && stringPayments == false { // Default case
		// Check for folder and open it, if it doesn't exist, Create one
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

		line := fmt.Sprintf("%d | %s\n", d.ExsOffset, stringData)
		d.ExsOffset++
		file.WriteString(line)

		c.JSON(201, gin.H{
			"message": "message sent successfully",
		})
	}

}

func (d *Demo) Consumer(c *gin.Context) {
	topic := c.Query("topic")
	offsetString := c.Query("offset")

	file, err := os.Open("internals/files/test.log")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "There was some error opening the file",
		})
		return
	}
	// scanner := bufio.NewScanner(file)
	defer file.Close()

	d.mu.Lock()
	defer d.mu.Unlock()

	var offset int

	if offsetString == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "failed to convert from string to integer",
			})
			return
		}
	}

	if topic == "orders" { // For Orders
		var OrdersPath string = "internals/files/folders/orders"
		OrdersText, _ := ConsumerReadFiles(c, OrdersPath, offset)
		c.JSON(200, gin.H{
			"message": OrdersText,
		})
		return

	} else if topic == "payments" { // For Payments
		var PaymentsPath string = "internals/files/folders/payments"
		PaymentsText, _ := ConsumerReadFiles(c, PaymentsPath, offset)
		c.JSON(200, gin.H{
			"message": PaymentsText,
		})
		return

	} else if topic == "" { // For Default
		var DefaultPath string = "internals/files/folders/default"
		DefaultText, _ := ConsumerReadFiles(c, DefaultPath, offset)
		c.JSON(200, gin.H{
			"message": DefaultText,
		})
		return
	} else {
		c.JSON(400, gin.H{
			"error": "no such URL key exists",
		})
		return
	}

	// if offsetString == "" {
	// 	// If offset is not present, scan through the entire file
	// 	for scanner.Scan() {
	// 		line := scanner.Text()
	// 		c.JSON(200, gin.H{
	// 			"message": line,
	// 		})
	// 	}
	// 	return
	// }

	// var cnt int = 0
	// // If offset is present, start scanning from the offset position of the file
	// for scanner.Scan() {
	// 	if cnt >= offset {
	// 		line := scanner.Text()
	// 		c.JSON(200, gin.H{
	// 			"message": line,
	// 		})
	// 	}
	// 	cnt++
	// }

}
