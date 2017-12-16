package tmsh

import (
	"fmt"
	"reflect"
	"testing"
)

var validCmds = []string{
	"show sys clock",
}

type TestSSHConnection struct {
	validCmd bool
	ret      []byte
}

func (c *TestSSHConnection) Send(cmd string) (int, error) {
	for _, vc := range validCmds {
		if vc == cmd {
			c.validCmd = true
			return 0, nil
		}
	}
	c.validCmd = false
	return 1, nil
}

func (c *TestSSHConnection) Recv(suffix string) ([]byte, error) {
	if !c.validCmd {
		return nil, fmt.Errorf("Syntax Error: unexpected argument \"foo\"")
	}
	return c.ret, nil
}

func (c *TestSSHConnection) Close() error {
	return nil
}

func TestExecuteCommand(t *testing.T) {
	bigip := &BigIP{
		host: "example.com",
		user: "admin.prd",
		sshconn: &TestSSHConnection{
			ret: []byte("Last login: Mon Jul  1 10:20:34 2017 from 192.0.2.10\nadmin.prd@(LB000)(cfg-sync In Sync)(/S1-green-P:Active)(/admin.prd)(tmos)#\nshow sys clock\n------------------------\nSys::Clock\n------------------------\nMon Jul 03 14:25:02 JST 2017\n\nn.prd)(tmos)# ve)(/admi\n"),
		},
	}

	expect := "------------------------\nSys::Clock\n------------------------\nMon Jul 03 14:25:02 JST 2017\n"
	actual, err := bigip.ExecuteCommand("show sys clock")
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %v\nwant %v", actual, expect)
	}

	expectErr := fmt.Errorf("Syntax Error: unexpected argument \"foo\"")
	ret, actualErr := bigip.ExecuteCommand("show foo")

	if ret != "" {
		t.Errorf("got %v\nwant %v", ret, "")
	}

	if !reflect.DeepEqual(actualErr, expectErr) {
		t.Errorf("got %v\nwant %v", actualErr, expectErr)
	}
}
