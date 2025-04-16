package ios

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/solywsh/go-forensic/utils"
	usbmuxd "github.com/solywsh/go-forensic/utils/iosutils/usbmuxd-device"
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/printer"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func handleDeviceList(cmd *cobra.Command, args []string) {
	usbHub := usbmuxd.NewUSBHub()
	deviceList, err := usbHub.DeviceList()
	if err != nil {
		log.Error(err)
	}
	if !osx.IsTTY() {
		for _, device := range deviceList {
			fmt.Printf("Id: %s\tConnectionSpeed: %s, ConnectionType: %s\n", device.SerialNumber, formatSpeed(device.ConnectionSpeed), device.ConnectionType)
		}
		return
	}
	columns := []table.Column{
		{Title: "Id", Width: utils.Min(50, utils.Max(6, len("00000000-0000000000000000")))},
		{Title: "ConnectionSpeed", Width: utils.Min(50, utils.Max(6, len("ConnectionSpeed")))},
		{Title: "ConnectionType", Width: utils.Min(50, utils.Max(6, len("ConnectionType")))},
	}
	rows := make([]table.Row, 0, len(deviceList))
	for _, r := range deviceList {
		rows = append(rows, table.Row{r.SerialNumber, formatSpeed(r.ConnectionSpeed), string(r.ConnectionType)})
	}
	tb := printer.NewTable(context.Background())
	if tableHeight > 0 {
		tb.SetHeight(tableHeight)
	}
	tb.SetColumns(columns).SetRows(rows)
	tb.Run()
	tb.Wait()
}

func formatSpeed(speed int) string {
	cast.ToString(speed / 1000000)
	if speed < 1000 {
		return fmt.Sprintf("%d b/s", speed)
	} else if speed < 1000000 {
		return fmt.Sprintf("%d Kb/s", speed/1000)
	} else if speed < 1000000000 {
		return fmt.Sprintf("%d Mb/s", speed/1000000)
	} else {
		return fmt.Sprintf("%d Gb/s", speed/1000000000)
	}
}
