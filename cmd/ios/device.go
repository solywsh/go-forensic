package ios

import (
	"github.com/solywsh/go-forensic/ios"
	"github.com/spf13/cobra"
)

var (
	tableHeight int
)

var (
	deviceListCmd = &cobra.Command{
		Use:   "list",
		Short: "list all connected iOS devices",
		Run: func(cmd *cobra.Command, args []string) {
			deviceHelper := ios.NewDeviceHelper(ios.WithTableHeight(tableHeight))
			err := deviceHelper.HandleDeviceList()
			if err != nil {
				log.Debug(err)
			}
		},
	}
	deviceCmd = &cobra.Command{
		Use:   "device",
		Short: "handling iOS devices",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	deviceListCmd.Flags().IntVarP(&tableHeight, "table-height", "t", 10, "the height of the table")
	deviceCmd.AddCommand(deviceListCmd)
	SystemIOSCmd.AddCommand(deviceCmd)
}
