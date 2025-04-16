package ios

import (
	"github.com/solywsh/go-forensic/utils/logger"
	"github.com/spf13/cobra"
)

var (
	log = logger.NewLogger()
)

var SystemIOSCmd = &cobra.Command{
	Use:   "ios",
	Short: "go-forensic processes commands related to the iOS system",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
}
