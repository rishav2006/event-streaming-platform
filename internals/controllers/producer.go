package controllers

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

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
				"partition": PartitionNumFinder(file.Name(), "orders"),
				"offset":    d.OffsetSliceOrders[0],
				"message":   "message successfully sent to orders catalog",
			})
			file.Close()

		} else if d.CounterOrder%3 == 1 {
			file, _ := CheckFolder(orders[1], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSliceOrders[1], stringData)
			file.WriteString(line)
			fmt.Println(line, 1)
			d.OffsetSliceOrders[1]++
			c.JSON(201, gin.H{
				"partition": PartitionNumFinder(file.Name(), "orders"),
				"offset":    d.OffsetSliceOrders[1],
				"message":   "message successfully sent to orders catalog",
			})
			file.Close()

		} else if d.CounterOrder%3 == 2 {
			file, _ := CheckFolder(orders[2], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSliceOrders[2], stringData)
			file.WriteString(line)
			fmt.Println(line, 2)
			d.OffsetSliceOrders[2]++
			c.JSON(201, gin.H{
				"partition": PartitionNumFinder(file.Name(), "orders"),
				"offset":    d.OffsetSliceOrders[2],
				"message":   "message successfully sent to orders catalog",
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
				"partition": PartitionNumFinder(file.Name(), "payments"),
				"offset":    d.OffsetSlicePayments[1],
				"message":   "message successfully sent to payments catalog",
			})

			file.Close()
		} else if d.CounterPayment%2 == 0 {
			file, _ := CheckFolder(payments[0], c)

			line := fmt.Sprintf("%d | %s\n", d.OffsetSlicePayments[0], stringData)
			file.WriteString(line)
			d.OffsetSlicePayments[0]++
			c.JSON(201, gin.H{
				"partition": PartitionNumFinder(file.Name(), "payments"),
				"offset":    d.OffsetSlicePayments[0],
				"message":   "message successfully sent to payments catalog",
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
			"partition": PartitionNumFinder(file.Name(), "default"),
			"offset":    d.ExsOffset,
			"message":   "message successfully sent to default catalog",
		})
	}

}
