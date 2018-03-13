package tmsh

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// BigIP is a struct for session state
type BigIP struct {
	host    string
	user    string
	sshconn SSH
}

// NewKeySession is NewSession plus key handling
func NewKeySession(host, port, user string, key []byte) (*BigIP, error) {
       return GenSession(host, port, user, "", key)
}

// NewSession sets up new SSH session to BIG-IP TMSH
func NewSession(host, port, user, password string) (*BigIP, error) {
	return GenSession(host,port,user,password,[]byte{})
}

// GenSession handles either Password or SSH Key based..
func GenSession(host, port, user, password string, key []byte) (*BigIP, error) {
    sshconn, err := newSSHConnection(host+":"+port, user, password, key)
	
	if err != nil {
		return nil, err
	}

	ret, _ := sshconn.Recv("# ")
	if !strings.Contains(string(ret), "(tmos)") {
		_, err := sshconn.Send("tmsh")
		if err != nil {
			return nil, err
		}

		ret, err = sshconn.Recv("# ")
		if err != nil {
			return nil, err
		}
	}

	bigip := &BigIP{
		host:    host,
		user:    user,
		sshconn: sshconn,
	}

	// Suppress pager output
	if _, err := bigip.ExecuteCommand("modify cli preference pager disabled display-threshold 0"); err != nil {
		return nil, err
	}

	return bigip, nil
}

// ExecuteCommand is used to execute any TMSH commands
func (bigip *BigIP) ExecuteCommand(cmd string) (string, error) {
	promptSuffix := "# "

	_, err := bigip.sshconn.Send(cmd)
	if err != nil {
		return "", err
	}

	results, err := bigip.sshconn.Recv(promptSuffix)
	if err != nil {
		return "", err
	}

	results = removeCarriageReturn(results)
	reader := bytes.NewReader(results)
	scanner := bufio.NewScanner(reader)

	var lines []string

	for scanner.Scan() {
		text := scanner.Text()
		line := removeSpaceAndBackspace(text)

		if strings.HasPrefix(line, "Last login:") ||
			strings.Contains(line, "(tmos)") ||
			strings.HasPrefix(line, cmd) {
			continue
		}

		lines = append(lines, text)
	}

	return strings.Join(lines, "\n"), nil
}

// Save is used to execute 'save /sys config' command
func (bigip *BigIP) Save() error {
	ret, err := bigip.ExecuteCommand("save /sys config current-partition")
	if err != nil {
		fmt.Errorf(ret)
	}

	if strings.Contains(ret, "Syntax Error: \"current-partition\" unknown property") {
		ret, err = bigip.ExecuteCommand("save /sys config")
		if err != nil {
			fmt.Errorf(ret)
		}
	}

	if strings.Contains(ret, "Error") {
		return fmt.Errorf(ret)
	}

	return nil
}

// Close is used to close SSH session
func (bigip *BigIP) Close() {
	bigip.sshconn.Close()
}
