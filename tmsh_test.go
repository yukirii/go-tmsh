package tmsh

import (
	"fmt"
	"reflect"
	"testing"
)

type TestSSHConn struct {
	validCmd bool
}

func (conn *TestSSHConn) Send(cmd string) (int, error) {
	conn.validCmd = (cmd == "show sys clock")
	return 0, nil
}

func (conn *TestSSHConn) Recv(suffix string) ([]byte, error) {
	var ret []byte
	var err error

	if conn.validCmd {
		ret = []byte("Last login: Mon Jul  1 10:20:34 2017 from 192.0.2.10\nadmin.prd@(LB000)(cfg-sync In Sync)(/S1-green-P:Active)(/admin.prd)(tmos)#\nshow sys clock\n------------------------\nSys::Clock\n------------------------\nMon Jul 03 14:25:02 JST 2017\n\nn.prd)(tmos)# ve)(/admi\n")
		err = nil
	} else {
		ret = nil
		err = fmt.Errorf("Syntax Error: unexpected argument \"foo\"")
	}

	return ret, err
}

func (conn *TestSSHConn) Close() error {
	return nil
}

func TestExecuteCommand(t *testing.T) {
	bigip := &BigIP{
		host:    "example.com",
		user:    "admin.prd",
		sshconn: &TestSSHConn{},
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
	_, actualErr := bigip.ExecuteCommand("show foo")

	if !reflect.DeepEqual(actualErr, expectErr) {
		t.Errorf("got %v\nwant %v", actualErr, expectErr)
	}
}
