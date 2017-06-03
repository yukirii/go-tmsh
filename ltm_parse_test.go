package tmsh

import (
	"reflect"
	"testing"
)

func TestParseShowLtmNode(t *testing.T) {
	str := `ltm node dev-web01.example.com {
    addr 192.0.2.1
    cur-sessions 0
    monitor-rule none
    monitor-status unchecked
    name dev-web01.example.com
    serverside.bits-in 0
    serverside.bits-out 0
    serverside.cur-conns 0
    serverside.max-conns 0
    serverside.pkts-in 0
    serverside.pkts-out 0
    serverside.tot-conns 0
    session-status enabled
    status.availability-state unknown
    status.enabled-state enabled
    status.status-reason Node address does not have service checking enabled
    tot-requests 0
}`

	node := ParseShowLtmNode(str)

	expect := &Node{
		Addr:          "192.0.2.1",
		Name:          "dev-web01.example.com",
		MonitorRule:   "none",
		MonitorStatus: "unchecked",
		EnabledState:  "enabled",
	}

	if !reflect.DeepEqual(node, expect) {
		t.Errorf("got %v\nwant %v", node, expect)
	}
}
