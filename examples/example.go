package main

import (
	"fmt"

	"github.com/shiftky/go-tmsh"
)

func main() {
	// Connect to BIG-IP
	// Create New SSH Session
	bigip, err := tmsh.NewSession("lb01.example.com", "22", "admin", "secret")
	if err != nil {
		fmt.Println(err)
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

	// Create a pool, and add pool members.
	bigip.CreatePool("P_www.example.com_80")
	bigip.AddPoolMember("P_www.example.com_80", "web01.example.com", "http", 80)
	bigip.AddPoolMember("P_www.example.com_80", "web02.example.com", "http", 80)

	// Create a virtual server.
	bigip.CreateVirtualServer("VS_www.example.com_80", "P_www.example.com_80", "203.0.113.1", "tcp", 80)

	// Delete a pool member.
	bigip.DeletePoolMember("P_www.example.com_80", "web02.example.com", 80)

	// Delete a virtual server.
	bigip.DeleteVirtualServer("VS_www.example.com_80")

	// Delete a pool.
	bigip.DeletePool("P_www.example.com_80")

	// Delete nodes.
	bigip.DeleteNode("web01.example.com")
	bigip.DeleteNode("web02.example.com")
}
