package tmsh

import (
	"reflect"
	"testing"
)

func TestCurrentLine(t *testing.T) {
	str := `line1
line2
line3
line4`

	fm := NewFieldManager(str)

	actual := fm.CurrentLine()
	expect := "line1"

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %v\nwant %v", actual, expect)
	}

	fm.currentLineNum += 2

	actual = fm.CurrentLine()
	expect = "line3"

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %v\nwant %v", actual, expect)
	}
}

func TestAdvance(t *testing.T) {
	str := `line1
line2
line3
line4`

	fm := NewFieldManager(str)
	fm.Advance()
	fm.Advance()
	fm.Advance()

	actual := fm.currentLineNum
	expect := 3

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %v\nwant %v", actual, expect)
	}
}

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

func TestParseShowLtmPool(t *testing.T) {
	//# show ltm pool api.example.com_8080 members field-fmt
	str := `ltm pool api.example.com_8080 {
    active-member-cnt 2
    connq-all.age-edm 0
    connq-all.age-ema 0
    connq-all.age-head 0
    connq-all.age-max 0
    connq-all.depth 0
    connq-all.serviced 0
    connq.age-edm 0
    connq.age-ema 0
    connq.age-head 0
    connq.age-max 0
    connq.depth 0
    connq.serviced 0
    cur-sessions 0
    members {
        api01.example.com:8080 {
            addr 192.0.2.1
            connq.age-edm 0
            connq.age-ema 0
            connq.age-head 0
            connq.age-max 0
            connq.depth 0
            connq.serviced 0
            cur-sessions 0
            monitor-rule /Common/tcp (pool monitor)
            monitor-status up
            node-name api01.example.com
            pool-name api.example.com_8080
            port 8080
            serverside.bits-in 36.2K
            serverside.bits-out 87.9K
            serverside.cur-conns 0
            serverside.max-conns 3
            serverside.pkts-in 20
            serverside.pkts-out 20
            serverside.tot-conns 3
            session-status enabled
            status.availability-state available
            status.enabled-state enabled
            status.status-reason Pool member is available
            tot-requests 0
        }
        api02.example.com:8080 {
            addr 192.0.2.2
            connq.age-edm 0
            connq.age-ema 0
            connq.age-head 0
            connq.age-max 0
            connq.depth 0
            connq.serviced 0
            cur-sessions 0
            monitor-rule none
            monitor-status unchecked
            node-name api02.example.com
            pool-name api.example.com_8080
            port 8080
            serverside.bits-in 7.8M
            serverside.bits-out 44.5M
            serverside.cur-conns 0
            serverside.max-conns 42
            serverside.pkts-in 9.0K
            serverside.pkts-out 7.8K
            serverside.tot-conns 1.4K
            session-status user-disabled
            status.availability-state unknown
            status.enabled-state disabled
            status.status-reason Pool member does not have service checking enabled
            tot-requests 0
        }
    }
    min-active-members 0
    monitor-rule /Common/tcp
    name api.example.com_8080
    serverside.bits-in 7.8M
    serverside.bits-out 44.5M
    serverside.cur-conns 0
    serverside.max-conns 45
    serverside.pkts-in 9.0K
    serverside.pkts-out 7.8K
    serverside.tot-conns 1.4K
    status.availability-state available
    status.enabled-state enabled
    status.status-reason The pool is available
    tot-requests 0
}`

	pool := ParseShowLtmPool(str)

	poolMembers := []PoolMember{
		PoolMember{
			Name:              "api01.example.com",
			Addr:              "192.0.2.1",
			Port:              8080,
			MonitorRule:       "/Common/tcp (pool monitor)",
			MonitorStatus:     "up",
			EnabledState:      "enabled",
			AvailabilityState: "available",
			StatusReason:      "Pool member is available",
		},
		PoolMember{
			Name:              "api02.example.com",
			Addr:              "192.0.2.2",
			Port:              8080,
			MonitorRule:       "none",
			MonitorStatus:     "unchecked",
			EnabledState:      "disabled",
			AvailabilityState: "unknown",
			StatusReason:      "Pool member does not have service checking enabled",
		},
	}

	expect := &Pool{
		ActiveMemberCount: 2,
		MonitorRule:       "/Common/tcp",
		Name:              "api.example.com_8080",
		AvailabilityState: "available",
		EnabledState:      "enabled",
		StatusReason:      "The pool is available",
		PoolMembers:       poolMembers,
	}

	if !reflect.DeepEqual(pool, expect) {
		t.Errorf("\ngot %v\nwant %v", pool, expect)
	}
}

func TestParseListLtmVirtual(t *testing.T) {
	//# list ltm virtual api.example.com_80
	str := `ltm virtual api.example.com_80 {
    destination 203.0.113.1:http
    ip-protocol tcp
    mask 255.255.255.255
    partition partition1
    pool api.example.com_80
    profiles {
        /Common/tcp { }
    }
    source 0.0.0.0/0
    vs-index 1234
}`

	vs := ParseListLtmVirtual(str)

	expect := &VirtualServer{
		Destination: "203.0.113.1:http",
		IpProtocol:  "tcp",
		Mask:        "255.255.255.255",
		Partition:   "partition1",
		Pool:        "api.example.com_80",
	}

	if !reflect.DeepEqual(vs, expect) {
		t.Errorf("\ngot %v\nwant %v", vs, expect)
	}
}
