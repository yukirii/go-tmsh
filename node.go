package tmsh

import (
	"fmt"
	"strings"
)

type Node struct {
	Name          string `ltm:"name"`
	Addr          string `ltm:"addr"`
	MonitorRule   string `ltm:"monitor-rule"`
	MonitorStatus string `ltm:"monitor-status"`
	EnabledState  string `ltm:"status.enabled-state"`
}

func (bigip *BigIP) GetAllNodes() ([]Node, error) {
	ret, err := bigip.ExecuteCommand("show ltm node field-fmt")
	if err != nil {
		return nil, err
	}

	var nodes []Node
	for _, s := range splitLtmOutput(ret) {
		var n Node
		if err := Unmarshal(s, &n); err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}

	return nodes, nil
}

func (bigip *BigIP) GetNode(name string) (*Node, error) {
	ret, _ := bigip.ExecuteCommand("show ltm node " + name + " field-fmt")
	if strings.Contains(ret, "was not found.") {
		return nil, fmt.Errorf(ret)
	}

	var node Node
	if err := Unmarshal(ret, &node); err != nil {
		return nil, err
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
