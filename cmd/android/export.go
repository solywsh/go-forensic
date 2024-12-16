package android

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/solywsh/go-forensic/android"
	"github.com/solywsh/go-forensic/utils/printer"
	"github.com/spf13/cobra"
	"time"
)

var (
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "to export data from applications in the Android system",
		Run: func(cmd *cobra.Command, args []string) {
			ad := android.NewAndroid().SetOutput(output)
			defer func() {
				if len(specifyPaths) == 0 && len(keywords) == 0 {
					return
				}
				p := printer.NewSpinner().SetSpinner(spinner.Moon)
				p.Msg("done.").Run()
				time.Sleep(1 * time.Second)
				p.Quit()
			}()
			if len(specifyPaths) > 0 {
				err := ad.ExportBySpecify(specifyPaths...)
				if err != nil {
					log.Error(err)
					return
				}
			}
			if len(keywords) > 0 {
				err := ad.ExportAppDataByKeywords(keywords...)
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
	SystemAndroidCmd.AddCommand(exportCmd)
}
