package tmsh

import (
	"strings"
)

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
