package tmsh

import (
	"fmt"
	"strings"
)

type Node struct {
	Name          string
	Addr          string
	MonitorRule   string
	MonitorStatus string
	EnabledState  string
}

type Pool struct {
	ActiveMemberCount int
	Name              string
	MonitorRule       string
	AvailabilityState string
	EnabledState      string
	StatusReason      string
	PoolMembers       []PoolMember
}

type PoolMember struct {
	Name              string
	Addr              string
	Port              int
	MonitorRule       string
	MonitorStatus     string
	EnabledState      string
	AvailabilityState string
	StatusReason      string
}

func (bigip *BigIP) GetNode(name string) (*Node, error) {
	ret, _ := bigip.ExecuteCommand("show ltm node " + name + " field-fmt")
	if strings.Contains(ret, "was not found.") {
		return nil, fmt.Errorf(ret)
	}

	node := ParseShowLtmNode(ret)

	return node, nil
}

func (bigip *BigIP) CreateNode(name, ipaddr string) error {
	ret, _ := bigip.ExecuteCommand("create ltm node " + name + " address " + ipaddr)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) DeleteNode(name string) error {
	ret, _ := bigip.ExecuteCommand("delete ltm node " + name)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) EnableNode(name string) error {
	ret, _ := bigip.ExecuteCommand("modify ltm node " + name + " session user-enabled")
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) DisableNode(name string) error {
	ret, _ := bigip.ExecuteCommand("modify ltm node " + name + " session user-disabled")
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}
