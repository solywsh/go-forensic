package android

import (
	"github.com/solywsh/go-forensic/android"
	"github.com/spf13/cobra"
)

var (
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "to export data from applications in the Android system",
		Run: func(cmd *cobra.Command, args []string) {
			ad := android.NewAndroid()
			err := ad.SetOutput(output).ExportByKeywords(keywords...)
			if err != nil {
				log.Error(err)
				return
			}
		},
	}
)

func init() {
	exportCmd.Flags().StringVarP(&output, "output", "o", "temp", "the output directory for the exported data")
	exportCmd.Flags().StringSliceVarP(&keywords, "keyword", "k", nil, "the keyword to search for in the application data")
	SystemAndroidCmd.AddCommand(exportCmd)
}
