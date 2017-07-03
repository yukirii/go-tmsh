package tmsh

import (
	"bytes"
	"io"

	"golang.org/x/crypto/ssh"
)

type SSH interface {
	Send(cmd string) (int, error)
	Recv(suffix string) ([]byte, error)
	Close() error
}

type SSHConn struct {
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
	stderr  io.Reader
}

type KeyboardInteractive map[string]string

func (ki KeyboardInteractive) Challenge(user, instruction string, questions []string, echos []bool) ([]string, error) {
	var answers []string

	for _, q := range questions {
		answers = append(answers, ki[q])
	}

	return answers, nil
}

func NewSSHConnection(addr, user, password string) (SSH, error) {
	session, err := newSSHSession(addr, user, password)
	if err != nil {
		return nil, err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.ECHOCTL:       0,
		ssh.TTY_OP_ISPEED: 9600,
		ssh.TTY_OP_OSPEED: 9600,
	}

	err = session.RequestPty("xterm", 80, 24, modes)
	if err != nil {
		return nil, err
	}

	err = session.Shell()
	if err != nil {
		return nil, err
	}

	return &SSHConn{
		session: session,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
	}, nil
}

func newSSHSession(addr, user, password string) (*ssh.Session, error) {
	answers := KeyboardInteractive(map[string]string{
		"Password: ": password,
	})

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.KeyboardInteractive(
				answers.Challenge,
			),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (conn *SSHConn) Send(cmd string) (int, error) {
	return conn.stdin.Write([]byte(cmd + "\n"))
}

func (conn *SSHConn) Recv(suffix string) ([]byte, error) {
	var result bytes.Buffer
	buff := make([]byte, 65535)
	for {
		n, err := conn.stdout.Read(buff)
		if err != io.EOF && err != nil {
			return nil, err
		}
		result.Write(buff[:n])
		if err == io.EOF || bytes.HasSuffix(buff[:n], []byte(suffix)) {
			break
		}
	}
	return result.Bytes(), nil
}

func (conn *SSHConn) Close() error {
	return conn.session.Close()
}
