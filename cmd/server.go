package main

import (
	"fmt"
	ras "tp3/agt"
)

func main() {
	server := ras.NewServerRestAgent(":8080")
	//fmt.Println(bltAgt)
	server.Start()
	fmt.Scanln()
}
