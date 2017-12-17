package tmsh

import (
	"reflect"
	"testing"

	"github.com/k0kubun/pp"
)

func TestGetNode(t *testing.T) {
	//show ltm node dev-web01.example.com field-fmt
	retStr := `ltm node dev-web01.example.com {
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

	expect := &Node{
		Addr:          "192.0.2.1",
		Name:          "dev-web01.example.com",
		MonitorRule:   "none",
		MonitorStatus: "unchecked",
		EnabledState:  "enabled",
	}

	bigip := &BigIP{
		host: "example.com",
		user: "admin",
		sshconn: &TestSSHConnection{
			ret: []byte(retStr),
		},
	}

	node, err := bigip.GetNode("dev-web01.example.com")
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(node, expect) {
		t.Errorf("got :" + pp.Sprint(node))
		t.Errorf("want :" + pp.Sprint(expect))
	}
}

func TestGetPool(t *testing.T) {
	// show ltm pool api.example.com_8080 members field-fmt
	retStr := `ltm pool api.example.com_8080 {
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

	expect := &Pool{
		ActiveMemberCount: 2,
		MonitorRule:       "/Common/tcp",
		Name:              "api.example.com_8080",
		AvailabilityState: "available",
		EnabledState:      "enabled",
		StatusReason:      "The pool is available",
		PoolMembers: []PoolMember{
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
		},
	}

	bigip := &BigIP{
		host: "example.com",
		user: "admin",
		sshconn: &TestSSHConnection{
			ret: []byte(retStr),
		},
	}

	pool, err := bigip.GetPool("api.example.com_8080")
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(pool, expect) {
		t.Errorf("got :" + pp.Sprint(pool))
		t.Errorf("want :" + pp.Sprint(expect))
	}
}

func TestGetVirtual(t *testing.T) {
	// list ltm virtual api.example.com_443
	retStr := `ltm virtual api.example.com_443 {
	destination 203.0.113.1:https
	ip-protocol tcp
	mask 255.255.255.255
	partition partition1
	pool api.example.com_443
	profiles {
		/Common/tcp {
			context all
		}
		wildcard.example.com {
			context clientside
		}
	}
	source 0.0.0.0/0
	vs-index 1234
}`

	expect := &VirtualServer{
		Name:        "api.example.com_443",
		Destination: "203.0.113.1:https",
		IpProtocol:  "tcp",
		Mask:        "255.255.255.255",
		Partition:   "partition1",
		Pool:        "api.example.com_443",
		Profiles: map[string]Profile{
			"/Common/tcp":          Profile{Context: "all"},
			"wildcard.example.com": Profile{Context: "clientside"},
		},
	}

	bigip := &BigIP{
		host: "example.com",
		user: "admin",
		sshconn: &TestSSHConnection{
			ret: []byte(retStr),
		},
	}

	vs, err := bigip.GetVirtualServer("api.example.com_443")
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(vs, expect) {
		t.Errorf("got :" + pp.Sprint(vs))
		t.Errorf("want :" + pp.Sprint(expect))
	}
}

func TestGetAllNodes(t *testing.T) {
	// show ltm node field-fmt
	retStr := `ltm node dev-web01.example.com {
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
}
ltm node dev-web02.example.com {
    addr 192.0.2.2
    cur-sessions 0
    monitor-rule none
    monitor-status unchecked
    name dev-web02.example.com
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

	expect := []Node{
		Node{
			Addr:          "192.0.2.1",
			Name:          "dev-web01.example.com",
			MonitorRule:   "none",
			MonitorStatus: "unchecked",
			EnabledState:  "enabled",
		},
		Node{
			Addr:          "192.0.2.2",
			Name:          "dev-web02.example.com",
			MonitorRule:   "none",
			MonitorStatus: "unchecked",
			EnabledState:  "enabled",
		},
	}

	bigip := &BigIP{
		host: "example.com",
		user: "admin",
		sshconn: &TestSSHConnection{
			ret: []byte(retStr),
		},
	}

	nodes, err := bigip.GetAllNodes()
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(nodes, expect) {
		t.Errorf("got :" + pp.Sprint(nodes))
		t.Errorf("want :" + pp.Sprint(expect))
	}
}

func TestGetAllPools(t *testing.T) {
	// show ltm pool members field-fmt
	retStr := `ltm pool api.example.com_8080 {
    active-member-cnt 1
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
            serverside.bits-in 0
            serverside.bits-out 0
            serverside.cur-conns 0
            serverside.max-conns 0
            serverside.pkts-in 0
            serverside.pkts-out 0
            serverside.tot-conns 0
            session-status enabled
            status.availability-state available
            status.enabled-state enabled
            status.status-reason Pool member is available
            tot-requests 0
        }
    }
    min-active-members 0
    monitor-rule /Common/tcp
    name api.example.com_8080
    serverside.bits-in 0
    serverside.bits-out 0
    serverside.cur-conns 0
    serverside.max-conns 0
    serverside.pkts-in 0
    serverside.pkts-out 0
    serverside.tot-conns 0
    status.availability-state available
    status.enabled-state enabled
    status.status-reason The pool is available
    tot-requests 0
}
ltm pool web.example.com_80 {
    active-member-cnt 1
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
        web.example.com:80 {
            addr 192.0.2.2
            connq.age-edm 0
            connq.age-ema 0
            connq.age-head 0
            connq.age-max 0
            connq.depth 0
            connq.serviced 0
            cur-sessions 0
            monitor-rule /Common/tcp (pool monitor)
            monitor-status up
            node-name web01.example.com
            pool-name web.example.com_80
            port 80
            serverside.bits-in 0
            serverside.bits-out 0
            serverside.cur-conns 0
            serverside.max-conns 0
            serverside.pkts-in 0
            serverside.pkts-out 0
            serverside.tot-conns 0
            session-status enabled
            status.availability-state available
            status.enabled-state enabled
            status.status-reason Pool member is available
            tot-requests 0
        }
    }
    min-active-members 0
    monitor-rule /Common/tcp
    name web.example.com_80
    serverside.bits-in 0
    serverside.bits-out 0
    serverside.cur-conns 0
    serverside.max-conns 0
    serverside.pkts-in 0
    serverside.pkts-out 0
    serverside.tot-conns 0
    status.availability-state available
    status.enabled-state enabled
    status.status-reason The pool is available
    tot-requests 0
}`

	expect := []Pool{
		Pool{
			ActiveMemberCount: 1,
			MonitorRule:       "/Common/tcp",
			Name:              "api.example.com_8080",
			AvailabilityState: "available",
			EnabledState:      "enabled",
			StatusReason:      "The pool is available",
			PoolMembers: []PoolMember{
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
			},
		},
		Pool{
			ActiveMemberCount: 1,
			MonitorRule:       "/Common/tcp",
			Name:              "web.example.com_80",
			AvailabilityState: "available",
			EnabledState:      "enabled",
			StatusReason:      "The pool is available",
			PoolMembers: []PoolMember{
				PoolMember{
					Name:              "web01.example.com",
					Addr:              "192.0.2.2",
					Port:              80,
					MonitorRule:       "/Common/tcp (pool monitor)",
					MonitorStatus:     "up",
					EnabledState:      "enabled",
					AvailabilityState: "available",
					StatusReason:      "Pool member is available",
				},
			},
		},
	}

	bigip := &BigIP{
		host: "example.com",
		user: "admin",
		sshconn: &TestSSHConnection{
			ret: []byte(retStr),
		},
	}

	pools, err := bigip.GetAllPools()
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(pools, expect) {
		t.Errorf("got :" + pp.Sprint(pools))
		t.Errorf("want :" + pp.Sprint(expect))
	}
}

func TestGetAllVirtualServers(t *testing.T) {
	// list ltm virtual all-properties
	retStr := `ltm virtual api.example.com_443 {
    address-status yes
    app-service none
    auth none
    auto-lasthop default
    bwc-policy none
    clone-pools none
    cmp-enabled yes
    connection-limit 0
    description none
    destination 203.0.113.1:https
    enabled
    fallback-persistence none
    fw-enforced-policy api.example.com_443
    fw-staged-policy none
    gtm-score 0
    ip-intelligence-policy none
    ip-protocol tcp
    last-hop-pool none
    mask 255.255.255.255
    metadata none
    mirror disabled
    mobile-app-tunnel disabled
    nat64 disabled
    partition partition1
    persist none
    policies none
    pool api.example.com_8080
    profiles {
        /Common/tcp {
            context all
        }
    }
    rate-class none
    rate-limit disabled
    rate-limit-dst-mask 0
    rate-limit-mode object
    rate-limit-src-mask 0
    related-rules none
    rules none
    security-log-profiles none
    source 0.0.0.0/0
    source-address-translation {
        pool /Common/partition1_snat
        type snat
    }
    source-port change
    syn-cookie-status not-activated
    traffic-classes none
    translate-address enabled
    translate-port enabled
    vlans none
    vlans-disabled
    vs-index 222
}
ltm virtual web.example.com_80 {
    address-status yes
    app-service none
    auth none
    auto-lasthop default
    bwc-policy none
    clone-pools none
    cmp-enabled yes
    connection-limit 0
    description none
    destination 203.0.113.2:http
    enabled
    fallback-persistence none
    fw-enforced-policy web.example.com_80
    fw-staged-policy none
    gtm-score 0
    ip-intelligence-policy none
    ip-protocol tcp
    last-hop-pool none
    mask 255.255.255.255
    metadata none
    mirror disabled
    mobile-app-tunnel disabled
    nat64 disabled
    partition partition1
    persist none
    policies none
    pool web.example.com_80
    profiles {
        /Common/tcp {
            context all
        }
    }
    rate-class none
    rate-limit disabled
    rate-limit-dst-mask 0
    rate-limit-mode object
    rate-limit-src-mask 0
    related-rules none
    rules none
    security-log-profiles none
    source 0.0.0.0/0
    source-address-translation {
        pool /Common/partition1_snat
        type snat
    }
    source-port change
    syn-cookie-status not-activated
    traffic-classes none
    translate-address enabled
    translate-port enabled
    vlans none
    vlans-disabled
    vs-index 222
}`

	expect := []VirtualServer{
		VirtualServer{
			Name:        "api.example.com_443",
			Destination: "203.0.113.1:https",
			IpProtocol:  "tcp",
			Mask:        "255.255.255.255",
			Partition:   "partition1",
			Pool:        "api.example.com_8080",
			Profiles: map[string]Profile{
				"/Common/tcp": Profile{
					Context: "all",
				},
			},
		},
		VirtualServer{
			Name:        "web.example.com_80",
			Destination: "203.0.113.2:http",
			IpProtocol:  "tcp",
			Mask:        "255.255.255.255",
			Partition:   "partition1",
			Pool:        "web.example.com_80",
			Profiles: map[string]Profile{
				"/Common/tcp": Profile{
					Context: "all",
				},
			},
		},
	}

	bigip := &BigIP{
		host: "example.com",
		user: "admin",
		sshconn: &TestSSHConnection{
			ret: []byte(retStr),
		},
	}

	vss, err := bigip.GetAllVirtualServers()
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(vss, expect) {
		t.Errorf("got :" + pp.Sprint(vss))
		t.Errorf("want :" + pp.Sprint(expect))
	}
}
