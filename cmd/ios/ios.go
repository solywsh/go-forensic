package ios

import (
	"github.com/solywsh/go-forensic/utils/logger"
	"github.com/spf13/cobra"
)

var (
	log = logger.NewLogger()
)

var (
	sshAdd       string
	sshPass      string
	keywords     []string
	specifyPaths []string
	output       string
)

var SystemIOSCmd = &cobra.Command{
	Use:   "ios",
	Short: "go-forensic processes commands related to the iOS system",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	SystemIOSCmd.PersistentFlags().StringVarP(&sshAdd, "addr", "a", "root@127.0.0.1:22", "the username and address of the iOS device")
	SystemIOSCmd.PersistentFlags().StringVarP(&sshPass, "pass", "p", "alpine", "the password of the iOS device")
}
