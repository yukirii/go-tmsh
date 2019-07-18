package tmsh

import (
	"fmt"
	"strconv"
	"strings"
)

// Pool contains information about each Pool
type Pool struct {
	ActiveMemberCount int          `ltm:"active-member-cnt"`
	Name              string       `ltm:"name"`
	MonitorRule       string       `ltm:"monitor-rule"`
	AvailabilityState string       `ltm:"status.availability-state"`
	EnabledState      string       `ltm:"status.enabled-state"`
	StatusReason      string       `ltm:"status.status-reason"`
	PoolMembers       []PoolMember `ltm:"members"`
}

// Pool contains information about each pool member in a pool
type PoolMember struct {
	Name              string `ltm:"node-name"`
	Addr              string `ltm:"addr"`
	Port              int    `ltm:"port"`
	MonitorRule       string `ltm:"monitor-rule"`
	MonitorStatus     string `ltm:"monitor-status"`
	EnabledState      string `ltm:"status.enabled-state"`
	AvailabilityState string `ltm:"status.availability-state"`
	StatusReason      string `ltm:"status.status-reason"`
}

// GetAllPools returns a list of all pools
func (bigip *BigIP) GetAllPools() ([]Pool, error) {
	ret, err := bigip.ExecuteCommand("show ltm pool members field-fmt")
	if err != nil {
		return nil, err
	}

	var pools []Pool
	for _, s := range splitLtmOutput(ret) {
		var p Pool
		if err := Unmarshal(s, &p); err != nil {
			return nil, err
		}
		pools = append(pools, p)
	}

	return pools, nil
}

// GetPool gets a pool by name. Return nil if the pool does not exist.
func (bigip *BigIP) GetPool(name string) (*Pool, error) {
	ret, _ := bigip.ExecuteCommand("show ltm pool " + name + " members field-fmt")
	if strings.Contains(ret, "was not found.") {
		return nil, fmt.Errorf(ret)
	}

	var pool Pool
	if err := Unmarshal(ret, &pool); err != nil {
		return nil, err
	}

	return &pool, nil
}

// CreatePool creates a new pool.
func (bigip *BigIP) CreatePool(name string) error {
	ret, _ := bigip.ExecuteCommand("create ltm pool " + name)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

// DeletePool removes a pool.
func (bigip *BigIP) DeletePool(name string) error {
	ret, _ := bigip.ExecuteCommand("delete ltm pool " + name)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

// AddMonitorToPool adds a monitor to pool.
func (bigip *BigIP) AddMonitorToPool(poolName, monitorName string) error {
	ret, _ := bigip.ExecuteCommand("modify ltm pool " + poolName + " monitor '" + monitorName + "'")
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

// AddPoolMember adds a new pool member.
func (bigip *BigIP) AddPoolMember(poolName, nodeName, monitorName string, port int, priority int) error {
	member := nodeName + ":" + strconv.Itoa(port)
	cmd := "modify ltm pool " + poolName + " members add { " + member + " { priority-group " + strconv.Itoa(priority) + " } } monitor '" + monitorName + "'"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

// DeletePoolMember removes a pool member.
func (bigip *BigIP) DeletePoolMember(poolName, nodeName string, port int) error {
	member := nodeName + ":" + strconv.Itoa(port)
	cmd := "modify ltm pool " + poolName + " members delete { " + member + " }"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

// EnablePoolMember changes the status of pool member to enable.
func (bigip *BigIP) EnablePoolMember(poolName, nodeName string, port int) error {
	member := nodeName + ":" + strconv.Itoa(port)
	cmd := "modify ltm pool " + poolName + " members modify { " + member + " { session user-enabled } }"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

// DisablePoolMember changes the status of pool member to disable.
func (bigip *BigIP) DisablePoolMember(poolName, nodeName string, port int) error {
	member := nodeName + ":" + strconv.Itoa(port)
	cmd := "modify ltm pool " + poolName + " members modify { " + member + " { session user-disabled } }"
	ret, _ := bigip.ExecuteCommand(cmd)
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}
