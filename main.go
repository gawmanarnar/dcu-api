// main.go

package main

func main() {
	server := Server{}
	server.Init("admin", "admin256", "dcu")
	server.Run()
}
