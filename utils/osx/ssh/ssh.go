package ssh

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"strings"
)

type (
	Client struct {
		host           string
		port           string
		username       string
		password       string
		privateKeyPath string

		config *ssh.ClientConfig
		client *ssh.Client
	}
	Option func(*Client)
)

func NewClient(options ...Option) *Client {
	c := &Client{
		host:           "127.0.0.1",
		port:           "22",
		username:       "root",
		password:       "",
		privateKeyPath: "~/.ssh/id_rsa",
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func WithHost(host string) Option {
	return func(x *Client) {
		x.host = host
	}
}

func WithPort(port string) Option {
	return func(x *Client) {
		x.port = port
	}
}

func WithUsername(username string) Option {
	return func(x *Client) {
		x.username = username
	}
}

func WithPassword(password string) Option {
	return func(x *Client) {
		x.password = password
	}
}

func WithPrivateKeyPath(privateKeyPath string) Option {
	return func(x *Client) {
		x.privateKeyPath = privateKeyPath
	}
}

// C connect
func (s *Client) C() (*Client, error) {
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

func (s *Client) Close() error {
	if s.client == nil {
		return nil
	}
	return s.client.Close()
}

func (s *Client) MustC() *Client {
	c, _ := s.C()
	return c
}

func (s *Client) ExecuteCommand(command string) (string, error) {
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

func (s *Client) ExecuteCommandWithStream(command string, stdout, stderr io.Writer) error {
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

func (s *Client) DeleteFiles(path string) error {
	_, err := s.ExecuteCommand(fmt.Sprintf("rm -rf %s", path))
	if err != nil {
		return fmt.Errorf("failed to delete remote file: %v", err)
	}
	return nil
}

func (s *Client) SFTP() *SftpX {
	if s.client == nil {
		return nil
	}
	sftpClient, err := sftp.NewClient(s.client)
	if err != nil {
		return nil
	}
	return &SftpX{
		Client:     s,
		sftpClient: sftpClient,
	}
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
