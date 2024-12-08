package ios

import (
	"github.com/solywsh/go-forensic/ios"
	"github.com/spf13/cobra"
)

var (
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "To export data from applications in the iOS system.",
		Run: func(cmd *cobra.Command, args []string) {
			username, host, port, err := parseSSHAddress(sshAdd)
			if err != nil {
				log.Error(err)
				return
			}
			err = ios.NewIOS().
				SetHost(host).
				SetPort(port).
				SetUsername(username).
				SetPasswd(sshPass).
				SetOutput(output).ExportAppDataByKeywords(keywords...)
			if err != nil {
				log.Error(err)
				return
			}
		},
	}
)

func init() {
	exportCmd.Flags().StringSliceVarP(&keywords, "keyword", "k", nil, "The keyword to search for in the application data.")
	exportCmd.Flags().StringVarP(&output, "output", "o", "temp", "The output directory for the exported data.")
	RootCmd.AddCommand(exportCmd)
}
