package main

import (
	"fmt"
	"github.com/noams0/Mini-project-IA04/agt"
)

func main() {
	server := agt.NewServerRestAgent(":8080")
	server.Start()
	fmt.Scanln()
}
