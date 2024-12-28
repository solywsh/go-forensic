package ios

import (
	"github.com/solywsh/go-forensic/utils/logger"
	"github.com/solywsh/go-forensic/utils/osx"
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

type IOS struct {
	username       string
	passwd         string
	host           string
	port           string
	privateKeyPath string
	output         string

	sshClient  *osx.SShX
	sftpClient *osx.SftpX
	keywords   []string

	spinner       *printer.SpinnerX
	cleanPathOnce sync.Once
}

func NewIOS() *IOS {
	return &IOS{
		username:       "root",
		passwd:         "alpine",
		host:           "127.0.0.1",
		port:           "22",
		privateKeyPath: "",
		output:         "./temp",
	}
}

func (t *IOS) SetUsername(username string) *IOS {
	t.username = username
	return t
}

func (t *IOS) SetPasswd(passwd string) *IOS {
	t.passwd = passwd
	return t
}

func (t *IOS) SetHost(host string) *IOS {
	t.host = host
	return t
}

func (t *IOS) SetPort(port string) *IOS {
	t.port = port
	return t
}

func (t *IOS) SetPrivateKeyPath(privateKeyPath string) *IOS {
	t.privateKeyPath = privateKeyPath
	return t
}

func (t *IOS) SetOutput(output string) *IOS {
	t.output = output
	return t
}

func (t *IOS) ssh() *osx.SShX {
	return osx.NewSSH().SetHost(t.host).SetPort(t.port).SetUsername(t.username).SetPassword(t.passwd).SetPrivateKeyPath(t.privateKeyPath)
}

func (t *IOS) removeOutput() {
	t.cleanPathOnce.Do(func() {
		t.spinner.Msg("cleaning " + t.output)
		osx.RemoveAll(t.output)
	})
}
