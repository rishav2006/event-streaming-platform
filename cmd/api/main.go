package main

import (
	"fmt"

	"github.com/rishav2006/event-streaming/internals/routes"
)

func main() {
	r := routes.NewRouter()
	r.Run(":8080")
	fmt.Println("Server is running on PORT:8080")
}
