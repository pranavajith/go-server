package main

import "fmt"

func main() {
	server := newAPIServer(":8000")
	server.Run()
	fmt.Println("Heyo!")
}
