package osx

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"strings"
)

type (
	SShX struct {
		host           string
		port           string
		username       string
		password       string
		privateKeyPath string

		config *ssh.ClientConfig
		client *ssh.Client
	}
)

func NewSSH() *SShX {
	return &SShX{
		host:           "127.0.0.1",
		port:           "22",
		username:       "root",
		password:       "alpine",
		privateKeyPath: "~/.ssh/id_rsa",
	}
}

func (s *SShX) SetHost(host string) *SShX {
	s.host = host
	return s
}

func (s *SShX) SetPort(port string) *SShX {
	s.port = port
	return s
}

func (s *SShX) SetUsername(username string) *SShX {
	s.username = username
	return s
}

func (s *SShX) SetPassword(password string) *SShX {
	s.password = password
	return s
}

func (s *SShX) SetPrivateKeyPath(privateKeyPath string) *SShX {
	s.privateKeyPath = privateKeyPath
	return s
}

// C connect
func (s *SShX) C() (*SShX, error) {
	if s.username == "" || s.host == "" || s.port == "" {
		return nil, fmt.Errorf("username, host and port are required")
	}
	if s.password == "" && s.privateKeyPath == "" {
		return nil, fmt.Errorf("password or private key is required")
	}
	config := &ssh.ClientConfig{
		User:            s.username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if s.password != "" {
		config.Auth = append(config.Auth, ssh.Password(s.password))
	} else {
		key, err := os.ReadFile(s.privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %v", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %v", err)
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}
	client, err := ssh.Dial("tcp", s.host+":"+s.port, config)
	if err != nil {
		return nil, err
	}
	s.config = config
	s.client = client

	return s, nil
}

func (s *SShX) Close() error {
	if s.client == nil {
		return nil
	}
	return s.client.Close()
}

func (s *SShX) MustC() *SShX {
	c, _ := s.C()
	return c
}

func (s *SShX) ExecuteCommand(command string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}
	return stdoutBuf.String(), nil
}

func (s *SShX) ExecuteCommandWithStream(command string, stdout, stderr io.Writer) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	session.Stdout = stdout
	session.Stderr = stderr

	err = session.Run(command)
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	return nil
}

func (s *SShX) DeleteFiles(path string) error {
	_, err := s.ExecuteCommand(fmt.Sprintf("rm -rf %s", path))
	if err != nil {
		return fmt.Errorf("failed to delete remote file: %v", err)
	}
	return nil
}

func ParseSSHAddress(sshAddr string) (string, string, string, error) {
	parts := strings.Split(sshAddr, "@")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid SSH address format")
	}

	username := parts[0]
	address := parts[1]

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse address: %w", err)
	}

	return username, host, port, nil
}
