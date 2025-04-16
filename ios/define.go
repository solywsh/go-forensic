package ios

import (
	"github.com/solywsh/go-forensic/utils/logger"
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/osx/ssh"
	"github.com/solywsh/go-forensic/utils/printer"
	"sync"
)

const (
	ApplicationPath       = "/private/var/mobile/Containers/Data/Application"
	AppGroupPath          = "/private/var/mobile/Containers/Shared/AppGroup"
	ContainerMetaDataFile = ".com.apple.mobile_container_manager.metadata.plist"

	ActiveFileName = "active.tar"
)

var (
	focusPathList = []string{
		ApplicationPath,
		AppGroupPath,
	}
	log = logger.NewLogger()
)

type SSHelper struct {
	username       string
	passwd         string
	host           string
	port           string
	privateKeyPath string
	output         string

	sshClient  *ssh.Client
	sftpClient *ssh.SftpX
	keywords   []string

	spinner       *printer.SpinnerX
	cleanPathOnce sync.Once
}

type SSHOption func(*SSHelper)

func NewSSHelper(options ...SSHOption) *SSHelper {
	s := &SSHelper{
		username:       "root",
		passwd:         "alpine",
		host:           "127.0.0.1",
		port:           "22",
		privateKeyPath: "",
		output:         "./temp",
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func WithHost(host string) SSHOption {
	return func(x *SSHelper) {
		x.host = host
	}
}

func WithPort(port string) SSHOption {
	return func(x *SSHelper) {
		x.port = port
	}
}

func WithUsername(username string) SSHOption {
	return func(x *SSHelper) {
		x.username = username
	}
}

func WithPassword(password string) SSHOption {
	return func(x *SSHelper) {
		x.passwd = password
	}
}

func WithOutput(output string) SSHOption {
	return func(x *SSHelper) {
		x.output = output
	}
}

func WithPrivateKeyPath(privateKeyPath string) SSHOption {
	return func(x *SSHelper) {
		x.privateKeyPath = privateKeyPath
	}
}

func (t *SSHelper) _ssh() *ssh.Client {
	return ssh.NewClient(
		ssh.WithHost(t.host),
		ssh.WithPort(t.port),
		ssh.WithUsername(t.username),
		ssh.WithPassword(t.passwd),
		ssh.WithPrivateKeyPath(t.privateKeyPath),
	)
}

func (t *SSHelper) removeOutput() {
	t.cleanPathOnce.Do(func() {
		t.spinner.Msg("cleaning " + t.output)
		osx.RemoveAll(t.output)
	})
}
