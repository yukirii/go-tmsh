package tmsh

import (
	"fmt"
	"strconv"
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

func (bigip *BigIP) GetPool(name string) (*Pool, error) {
	ret, _ := bigip.ExecuteCommand("show ltm pool " + name + " members field-fmt")
	if strings.Contains(ret, "was not found.") {
		return nil, fmt.Errorf(ret)
	}
	pool := ParseShowLtmPool(ret)
	return pool, nil
}

func (bigip *BigIP) CreatePool(name string) error {
	ret, _ := bigip.ExecuteCommand("create ltm pool " + name)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) DeletePool(name string) error {
	ret, _ := bigip.ExecuteCommand("delete ltm pool " + name)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) AddMonitorToPool(poolName, monitorName string) error {
	ret, _ := bigip.ExecuteCommand("modify ltm pool " + poolName + " monitor '" + monitorName + "'")
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) AddPoolMember(poolName, nodeName, monitorName string, port int) error {
	member := nodeName + ":" + strconv.Itoa(port)
	cmd := "modify ltm pool " + poolName + " members add { " + member + " } monitor '" + monitorName + "'"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) DeletePoolMember(poolName, nodeName string, port int) error {
	member := nodeName + ":" + strconv.Itoa(port)
	cmd := "modify ltm pool " + poolName + " members delete { " + member + " }"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) EnablePoolMember(poolName, nodeName string, port int) error {
	member := nodeName + ":" + strconv.Itoa(port)
	cmd := "modify ltm pool " + poolName + " members modify { " + member + " { session user-enabled } }"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (bigip *BigIP) DisablePoolMember(poolName, nodeName string, port int) error {
	member := nodeName + ":" + strconv.Itoa(port)
	cmd := "modify ltm pool " + poolName + " members modify { " + member + " { session user-disabled } }"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}
