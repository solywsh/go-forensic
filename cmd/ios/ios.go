package ios

import (
	"fmt"
	"github.com/solywsh/go-forensic/utils/logger"
	"github.com/spf13/cobra"
	"net"
	"strings"
)

var (
	log = logger.NewLogger()
)

var (
	sshAdd   string
	sshPass  string
	keywords []string
	output   string
)

var RootCmd = &cobra.Command{
	Use:   "ios",
	Short: "go-forensic processes commands related to the iOS system.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&sshAdd, "addr", "a", "root@127.0.0.1:22", "The username and address of the iOS device.")
	RootCmd.PersistentFlags().StringVarP(&sshPass, "pass", "p", "alpine", "The password of the iOS device.")
}

func parseSSHAddress(sshAddr string) (string, string, string, error) {
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
