package ios

import (
	"context"
	"fmt"
	"github.com/solywsh/go-forensic/ios"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
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
			deviceHelper := ios.NewDeviceHelper()
			ctx, cancel := context.WithCancel(context.Background())
			err := deviceHelper.HandleProxy(ctx, deviceId, protocol, remotePort, localPort)
			if err != nil {
				log.Error(err)
			}
			// Create a channel to receive signal notifications.
			sigChan := make(chan os.Signal, 1)
			// Pass the specified signal notification to the sigChan channel.
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			// A channel to determine when to exit.
			done := make(chan bool, 1)
			// Start a goroutine to handle signals.
			go func() {
				// Block until a signal is received.
				<-sigChan
				cancel()
				done <- true
			}()
			fmt.Println("the program is running, press ctrl+c to exit.")
			<-done
			fmt.Println()
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
