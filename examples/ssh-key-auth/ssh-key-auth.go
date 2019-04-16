package main

import (
	"fmt"
	"io/ioutil"

	"github.com/yukirii/go-tmsh"
)

// SSH Key-Based Authentication
func main() {
	// Read private key file
	key, err := ioutil.ReadFile("/Users/user/.ssh/id_rsa")
	if err != nil {
		panic(err)
	}

	// Connect to BIG-IP
	// Create New SSH Session using SSH key
	bigip, err := tmsh.NewKeySession("lb01.example.com", "22", "admin", key)
	if err != nil {
		panic(err)
	}
	defer bigip.Close()

	// Create nodes.
	bigip.CreateNode("web01.example.com", "192.0.2.1")
	bigip.CreateNode("web02.example.com", "192.0.2.2")

	// Disable Node
	bigip.DisableNode("web02.example.com")

	// Get node.
	node1, _ := bigip.GetNode("web01.example.com")
	fmt.Println("Name: " + node1.Name)
	fmt.Println("Addr: " + node1.Addr)
	fmt.Println("State: " + node1.EnabledState)

	node2, _ := bigip.GetNode("web02.example.com")
	fmt.Println("Name: " + node2.Name)
	fmt.Println("Addr: " + node2.Addr)
	fmt.Println("State: " + node2.EnabledState)

	// Delete nodes.
	bigip.DeleteNode("web01.example.com")
	bigip.DeleteNode("web02.example.com")
}
