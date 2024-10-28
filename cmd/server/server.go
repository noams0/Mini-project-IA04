package main

import (
	"fmt"
	agt "tp3/agt"
)

func main() {
	server := agt.NewServerRestAgent(":8080")
	server.Start()
	fmt.Scanln()
}
