package main

import "fmt"

func main() {
	fmt.Println("Hello from notification server")
	server := NewServer()
	server.StartServer("1234")
}
