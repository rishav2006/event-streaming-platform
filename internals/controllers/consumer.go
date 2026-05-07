package controllers

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (d *Demo) Consumer(c *gin.Context) {
	topic := c.Query("topic")
	offsetString := c.Query("offset")
	groupString := c.Query("group")

	// Group Divisions - Arrays (Hardcoded values)
	var GroupDivArrayOrder = [][]int{{1, 1, 1}, {2, 1, 0}, {1, 2, 0}, {3, 0, 0}}
	var GroupDivArrayPayment = [][]int{{1, 1}, {0, 2}, {2, 0}}

	// Group Division Mapping - Maps
	var GroupOrderMap = map[string][]int{"A": GroupDivArrayOrder[0], "B": GroupDivArrayOrder[1], "C": GroupDivArrayOrder[2], "D": GroupDivArrayOrder[3]}

	var GroupPaymentMap = map[string][]int{"A": GroupDivArrayPayment[0], "B": GroupDivArrayPayment[1], "C": GroupDivArrayPayment[2]}

	file, err := os.Open("internals/files/test.log")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "There was some error opening the file",
		})
		return
	}

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
		if groupString == "" {
			OrdersText, _ := ConsumerReadFilesNoGroup(c, OrdersPath, offset, topic)
			c.JSON(200, OrdersText)
		} else {
			OrdersText, _ := ConsumerReadFiles(c, OrdersPath, offset, groupString, GroupOrderMap)
			c.JSON(200, OrdersText)
		}
		return

	} else if topic == "payments" { // For Payments
		var PaymentsPath string = "internals/files/folders/payments"
		if groupString == "" {
			PaymentsText, _ := ConsumerReadFilesNoGroup(c, PaymentsPath, offset, topic)
			c.JSON(200, PaymentsText)
		} else {
			PaymentsText, _ := ConsumerReadFiles(c, PaymentsPath, offset, groupString, GroupPaymentMap)
			c.JSON(200, PaymentsText)
		}

		return

	} else if topic == "" { // For Default
		var DefaultPath string = "internals/files/folders/default"
		if groupString == "" {
			DefaultText, _ := ConsumerReadFilesNoGroup(c, DefaultPath, offset, topic)
			c.JSON(200, DefaultText)
		} else {
			DefaultText, _ := ConsumerReadFiles(c, DefaultPath, offset, groupString, GroupPaymentMap)
			c.JSON(200, DefaultText)
		}
		return
	} else {
		c.JSON(400, gin.H{
			"error": "no such URL key exists",
		})
		return
	}

}
