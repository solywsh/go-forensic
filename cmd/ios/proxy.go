package ios

import (
	"context"
	"fmt"
	"github.com/solywsh/go-forensic/ios"
	"github.com/solywsh/go-forensic/utils"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var (
	deviceId   string
	protocol   string
	remotePort int
	localPort  int
)

var (
	proxyCmd = &cobra.Command{
		Use:   "proxy",
		Short: "use USB to proxy the specified TCP/UDP port on the iOS device",
		Run: func(cmd *cobra.Command, args []string) {
			if protocol != "tcp" && protocol != "udp" {
				log.Error("protocol must be tcp or udp")
				return
			}
			deviceHelper := ios.NewDeviceHelper()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)

			if !utils.IsPortAvailable(protocol, localPort) {
				log.Warn("port may be occupied", "protocol", protocol, "port", localPort)
			}

			done := make(chan struct{})
			err := deviceHelper.HandleProxy(ctx, deviceId, protocol, remotePort, localPort)
			if err != nil {
				log.Error(err)
			}
			go func() {
				<-sigChan
				cancel()
				close(done)
			}()
			fmt.Println("the proxy is running, press ctrl+c to exit.")
			<-done
		},
	}
)

func init() {
	proxyCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "specify device Id, default is the first device detected")
	proxyCmd.Flags().StringVarP(&protocol, "protocol", "p", "tcp", "proxy protocol type, tcp/udp")
	proxyCmd.Flags().IntVarP(&localPort, "local-port", "l", 2222, "proxy local port")
	proxyCmd.Flags().IntVarP(&remotePort, "remote-port", "r", 22, "proxy remote device port")
	SystemIOSCmd.AddCommand(proxyCmd)
}
