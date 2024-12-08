package cmd

import (
	"github.com/solywsh/go-forensic/cmd/ios"
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
}

func init() {
	rootCmd.AddCommand(ios.RootCmd)
}

func Execute() {
	rootCmd.Execute()
}
