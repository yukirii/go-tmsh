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

func NewSession(host, user, password string) (*BigIP, error) {
	sshconn, err := NewSSHConnection(host, user, password)
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

	// Remove carriage return
	removedResults := make([]byte, 0)
	for _, c := range results {
		if c != 13 {
			removedResults = append(removedResults, c)
		}
	}

	reader := bytes.NewReader(removedResults)
	scanner := bufio.NewScanner(reader)

	returnStr := ""

	for scanner.Scan() {
		text := scanner.Text()

		// Remove space + back space
		line := strings.Replace(text, " \b", "", -1)

		if strings.HasPrefix(line, "Last login:") ||
			strings.HasPrefix(line, promptPrefix) ||
			strings.HasPrefix(line, cmd) {
			continue
		}

		returnStr += text + "\n"
	}

	return returnStr, nil
}

func (bigip *BigIP) Save() {
	bigip.sshconn.Close()
}
