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
	EnabledReason string
}

func (bigip *BigIP) GetNode(name string) (*Node, error) {
	ret, _ := bigip.ExecuteCommand("show ltm node " + name + " field-fmt")
	if strings.Contains(ret, "was not found.") {
		return nil, fmt.Errorf(ret)
	}

	node := Node{}
	for _, line := range strings.Split(ret, "\n") {
		if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "}") {
			continue
		}

		line = strings.TrimSpace(line)
		columns := strings.SplitAfterN(line, " ", 2)

		if strings.HasPrefix(columns[0], "addr") {
			node.Addr = columns[1]
		} else if strings.HasPrefix(columns[0], "name") {
			node.Name = columns[1]
		} else if strings.HasPrefix(columns[0], "monitor-rule") {
			node.MonitorRule = columns[1]
		} else if strings.HasPrefix(columns[0], "monitor-status") {
			node.MonitorStatus = columns[1]
		} else if strings.HasPrefix(columns[0], "status.enabled-state") {
			node.EnabledState = columns[1]
		} else if strings.HasPrefix(columns[0], "status.status-reason") {
			node.EnabledReason = columns[1]
		}
	}

	return &node, nil
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
