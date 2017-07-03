package tmsh

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type BigIP struct {
	host    string
	user    string
	sshconn SSH
}

func NewSession(host, port, user, password string) (*BigIP, error) {
	sshconn, err := NewSSHConnection(host+":"+port, user, password)
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

	return &BigIP{
		host:    host,
		user:    user,
		sshconn: sshconn,
	}, nil
}

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

func (bigip *BigIP) Close() {
	bigip.sshconn.Close()
}
