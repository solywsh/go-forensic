package android

import (
	"github.com/solywsh/go-forensic/utils/logger"
	"github.com/spf13/cobra"
)

var (
	log = logger.NewLogger()
)

var (
	keywords     []string
	specifyPaths []string
	output       string
)

var SystemAndroidCmd = &cobra.Command{
	Use:   "android",
	Short: "go-forensic processes commands related to the Android system",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
