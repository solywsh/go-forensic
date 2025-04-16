package cmd

import (
	"github.com/solywsh/go-forensic/cmd/android"
	"github.com/solywsh/go-forensic/cmd/db"
	"github.com/solywsh/go-forensic/cmd/ios"
	"github.com/solywsh/go-forensic/constant"
	"github.com/solywsh/go-forensic/utils/logger"
	"github.com/solywsh/go-forensic/utils/printer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-forensic",
	Short: "go-forensic is a supplement to mobile phone forensics tools.",
	Long: printer.WrapText("go-forensic is a tool developed to "+
		"address the lack of integrated features in some forensic tools, "+
		"with the intention of improving the efficiency of mobile forensics work.",
		65),
	Run: func(cmd *cobra.Command, args []string) {
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if constant.GetDebug() {
			log := logger.NewDebugLogger()
			log.Debug("start with debug mode..")
		}
	},
}

func init() {
	rootCmd.AddCommand(db.SqliteCmd)
	rootCmd.AddCommand(ios.SystemIOSCmd)
	rootCmd.AddCommand(android.SystemAndroidCmd)
}

func Execute() {
	rootCmd.PersistentFlags().BoolVar(&constant.Debug, "debug", false, "show debug info")
	rootCmd.Execute()
}
