package ios

import (
	"github.com/solywsh/go-forensic/ios"
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/printer"
	"github.com/spf13/cobra"
	"time"
)

var (
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "to export data from applications in the iOS system",
		Run: func(cmd *cobra.Command, args []string) {
			username, host, port, err := osx.ParseSSHAddress(sshAdd)
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
			iosHelper := ios.NewIOS().
				SetHost(host).
				SetPort(port).
				SetUsername(username).
				SetPasswd(sshPass).
				SetOutput(output)
			if len(specifyPaths) > 0 {
				err = iosHelper.ExportBySpecify(specifyPaths...)
				if err != nil {
					log.Error(err)
					return
				}
			}
			if len(keywords) > 0 {
				err = iosHelper.ExportAppDataByKeywords(keywords...)
				if err != nil {
					log.Error(err)
					return
				}
			}
		},
	}
)

func init() {
	exportCmd.Flags().StringVarP(&output, "output", "o", "temp", "the output directory for the exported data")
	exportCmd.Flags().StringSliceVarP(&keywords, "keyword", "k", nil, "the keyword to search for in the application data")
	exportCmd.Flags().StringSliceVarP(&specifyPaths, "specify-path", "s", nil, "the path to the data to be exported")
	SystemIOSCmd.AddCommand(exportCmd)
}
