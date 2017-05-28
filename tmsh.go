package tmsh

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"golang.org/x/crypto/ssh"
)

type BigIP struct {
	host    string
	user    string
	sshconn *SSHConn
}

type SSHConn struct {
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
	stderr  io.Reader
}

func NewSession(host, port, user, password string) (*BigIP, error) {
	sshconn, err := NewSSHConnection(host+":"+port, user, password)
	if err != nil {
		return nil, err
	}

	return &BigIP{
		host:    host,
		user:    user,
		sshconn: sshconn,
	}, nil
}

func (bigip *BigIP) ExecuteCommand(cmd string) (string, error) {
	promptSuffix := "# "
	promptPrefix := bigip.user + "@"

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
			strings.HasPrefix(line, promptPrefix) ||
			strings.HasPrefix(line, cmd) {
			continue
		}

		lines = append(lines, text)
	}

	return strings.Join(lines, "\n"), nil
}

func (bigip *BigIP) Close() {
	bigip.sshconn.Close()
}

func (bigip *BigIP) Save() {
	bigip.Close()
}
