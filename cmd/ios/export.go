package ios

import (
	"context"
	"github.com/solywsh/go-forensic/ios"
	"github.com/solywsh/go-forensic/utils/osx/ssh"
	"github.com/solywsh/go-forensic/utils/printer"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"time"
)

var (
	sshAddr           string
	sshPass           string
	keywords          []string
	specifyPaths      []string
	output            string
	host              string
	port              string
	username          string
	sshPrivateKeyPath string

	usbProxy bool
)

var (
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "to export data from applications in the iOS system",
		Run:   handleExport,
	}
)

func processSSHInfo() (string, string, string, error) {
	_username, _host, _port, err := ssh.ParseSSHAddress(sshAddr)
	if err != nil {
		return _username, _host, _port, err
	}
	if username != "" {
		_username = username
	}
	if host != "" {
		_host = host
	}
	if port != "" {
		_port = port
	}
	return _username, _host, _port, nil
}

func handleExport(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if usbProxy {
		deviceHelper := ios.NewDeviceHelper()
		deviceNum := deviceHelper.GetDeviceNum()
		if deviceNum > 0 {
			err := deviceHelper.HandleProxy(ctx, deviceId, "tcp", remotePort, localPort)
			if err != nil {
				log.Error(err)
				return
			}
			host = "127.0.0.1"
			port = cast.ToString(localPort)
		}
	}
	username, host, port, err := processSSHInfo()
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		p := printer.NewSpinner()
		p.Msg("done.").Run()
		time.Sleep(1 * time.Second)
		p.Quit()
	}()
	iosHelper := ios.NewSSHelper(
		ios.WithHost(host),
		ios.WithPort(port),
		ios.WithUsername(username),
		ios.WithPassword(sshPass),
		ios.WithOutput(output),
		ios.WithPrivateKeyPath(sshPrivateKeyPath),
	)
	if len(specifyPaths) > 0 {
		err = iosHelper.ExportBySpecify(specifyPaths...)
		if err != nil {
			log.Error(err)
		}
	}
	if len(keywords) > 0 {
		err = iosHelper.ExportAppDataByKeywords(keywords...)
		if err != nil {
			log.Error(err)
		}
	}
}

func init() {
	exportCmd.Flags().BoolVarP(&usbProxy, "usb-proxy", "u", true, "if the number of devices detected through the USB connection > 0, then use the USB proxy")
	exportCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "specify device Id, default is the first device detected")
	exportCmd.Flags().IntVarP(&localPort, "local-port", "l", 2222, "ssh port used for local proxy")
	exportCmd.Flags().IntVarP(&remotePort, "remote-port", "r", 22, "ssh port of the remote device")
	exportCmd.Flags().StringVarP(&sshAddr, "addr", "a", "root@127.0.0.1:22", "the username and address of the iOS device")
	exportCmd.Flags().StringVarP(&sshPass, "pass", "p", "alpine", "the ssh password of the iOS device")
	exportCmd.Flags().StringVar(&sshPrivateKeyPath, "private-key", "", "the ssh private key path")
	exportCmd.Flags().StringVar(&host, "host", "", "the ssh host of the iOS device")
	exportCmd.Flags().StringVar(&port, "port", "", "the ssh port of the iOS device")
	exportCmd.Flags().StringVar(&username, "username", "", "the ssh username of the iOS device")
	exportCmd.Flags().StringVarP(&output, "output", "o", "temp", "the output directory for the exported data")
	exportCmd.Flags().StringSliceVarP(&keywords, "keyword", "k", nil, "the keyword to search for in the application data")
	exportCmd.Flags().StringSliceVarP(&specifyPaths, "specify-path", "s", nil, "the path to the data to be exported")
	SystemIOSCmd.AddCommand(exportCmd)
}
