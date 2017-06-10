package tmsh

import (
	"fmt"
	"strconv"
	"strings"
)

type FieldManager struct {
	lines          []string
	currentLineNum int
}

func NewFieldManager(str string) FieldManager {
	return FieldManager{
		lines:          strings.Split(str, "\n"),
		currentLineNum: 0,
	}
}

func (fm *FieldManager) CurrentLine() string {
	return fm.lines[fm.currentLineNum]
}

func (fm *FieldManager) Advance() error {
	if len(fm.lines) == fm.currentLineNum {
		fmt.Errorf("Index out of range")
	}
	fm.currentLineNum += 1
	return nil
}

func ParseListLtmVirtual(str string) *VirtualServer {
	vs := VirtualServer{}

	fm := NewFieldManager(str)

	for len(fm.lines) > fm.currentLineNum {
		line := strings.TrimSpace(fm.CurrentLine())
		columns := strings.SplitAfterN(line, " ", 2)

		switch {
		case strings.HasPrefix(columns[0], "destination"):
			vs.Destination = columns[1]
		case strings.HasPrefix(columns[0], "ip-protocol"):
			vs.IpProtocol = columns[1]
		case strings.HasPrefix(columns[0], "mask"):
			vs.Mask = columns[1]
		case strings.HasPrefix(columns[0], "partition"):
			vs.Partition = columns[1]
		case strings.HasPrefix(columns[0], "pool"):
			vs.Pool = columns[1]
		}

		fm.Advance()
	}

	return &vs
}

func ParsePoolMemberNodes(fm *FieldManager) []PoolMember {
	poolMembers := []PoolMember{}

	if strings.HasSuffix(fm.CurrentLine(), "members {") {
		fm.Advance()
	}

	for {
		// Parse a Node
		if strings.HasSuffix(fm.CurrentLine(), "{") {
			var member PoolMember

			for {
				line := strings.TrimSpace(fm.CurrentLine())
				columns := strings.SplitAfterN(line, " ", 2)

				switch {
				case strings.HasPrefix(columns[0], "addr"):
					member.Addr = columns[1]
				case strings.HasPrefix(columns[0], "node-name"):
					member.Name = columns[1]
				case strings.HasPrefix(columns[0], "port"):
					member.Port, _ = strconv.Atoi(columns[1])
				case strings.HasPrefix(columns[0], "monitor-rule"):
					member.MonitorRule = columns[1]
				case strings.HasPrefix(columns[0], "monitor-status"):
					member.MonitorStatus = columns[1]
				case strings.HasPrefix(columns[0], "status.enabled-state"):
					member.EnabledState = columns[1]
				case strings.HasPrefix(columns[0], "status.availability-state"):
					member.AvailabilityState = columns[1]
				case strings.HasPrefix(columns[0], "status.status-reason"):
					member.StatusReason = columns[1]
				}

				if strings.HasSuffix(fm.CurrentLine(), "}") {
					poolMembers = append(poolMembers, member)
					break
				}

				fm.Advance()
			}
		}

		fm.Advance()

		// Check for end of "members {}"
		if strings.HasSuffix(fm.CurrentLine(), "}") {
			break
		}
	}

	return poolMembers
}

func ParseShowLtmPool(str string) *Pool {
	pool := Pool{}

	fm := NewFieldManager(str)

	for len(fm.lines) > fm.currentLineNum {
		// Parse Pool Members
		if strings.HasSuffix(fm.CurrentLine(), "members {") {
			pool.PoolMembers = ParsePoolMemberNodes(&fm)
		}

		line := strings.TrimSpace(fm.CurrentLine())
		columns := strings.SplitAfterN(line, " ", 2)

		// Parse Pool
		switch {
		case strings.HasPrefix(columns[0], "active-member-cnt"):
			pool.ActiveMemberCount, _ = strconv.Atoi(columns[1])
		case strings.HasPrefix(columns[0], "monitor-rule"):
			pool.MonitorRule = columns[1]
		case strings.HasPrefix(columns[0], "name"):
			pool.Name = columns[1]
		case strings.HasPrefix(columns[0], "status.availability-state"):
			pool.AvailabilityState = columns[1]
		case strings.HasPrefix(columns[0], "status.enabled-state"):
			pool.EnabledState = columns[1]
		case strings.HasPrefix(columns[0], "status.status-reason"):
			pool.StatusReason = columns[1]
		}

		fm.Advance()
	}

	return &pool
}

func ParseShowLtmNode(str string) *Node {
	node := Node{}
	for _, line := range strings.Split(str, "\n") {
		if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "}") {
			continue
		}

		line = strings.TrimSpace(line)
		columns := strings.SplitAfterN(line, " ", 2)

		switch {
		case strings.HasPrefix(columns[0], "addr"):
			node.Addr = columns[1]
		case strings.HasPrefix(columns[0], "name"):
			node.Name = columns[1]
		case strings.HasPrefix(columns[0], "monitor-rule"):
			node.MonitorRule = columns[1]
		case strings.HasPrefix(columns[0], "monitor-status"):
			node.MonitorStatus = columns[1]
		case strings.HasPrefix(columns[0], "status.enabled-state"):
			node.EnabledState = columns[1]
		}
	}

	return &node
}
