package ios

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/solywsh/go-forensic/utils"
	usbmuxd "github.com/solywsh/go-forensic/utils/iosutils/usbmuxd-device"
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/printer"
	"runtime"
	"strings"
	"time"
)

type DeviceHelper struct {
	tableHeight int
}

type DeviceOption func(*DeviceHelper)

func NewDeviceHelper(options ...DeviceOption) *DeviceHelper {
	d := &DeviceHelper{
		tableHeight: 10,
	}
	for _, option := range options {
		option(d)
	}
	return d
}

func WithTableHeight(height int) DeviceOption {
	return func(d *DeviceHelper) {
		d.tableHeight = height
	}
}

func (d *DeviceHelper) GetDeviceNum() int {
	usbHub := usbmuxd.NewUsbDriver()
	deviceList, err := usbHub.DeviceList()
	if err != nil {
		return 0
	}
	return len(deviceList)
}

func (d *DeviceHelper) HandleDeviceList() error {
	usbHub := usbmuxd.NewUsbDriver()
	deviceList, err := usbHub.DeviceList()
	if err != nil {
		return err
	}
	if len(deviceList) == 0 {
		return errors.New("no device found")
	}
	if !osx.IsTTY() {
		for _, device := range deviceList {
			fmt.Printf("Id: %s\tConnectionSpeed: %s, ConnectionType: %s\n", device.SerialNumber, utils.FormatNetSpeed(device.ConnectionSpeed), device.ConnectionType)
		}
		return nil
	}
	var columns []table.Column
	rows := make([]table.Row, 0, len(deviceList))

	if runtime.GOOS == "windows" {
		idLength := 6
		for _, r := range deviceList {
			idLength = utils.Max(idLength, len(r.SerialNumber))
			rows = append(rows, table.Row{r.SerialNumber, string(r.ConnectionType)})
		}
		columns = []table.Column{
			{Title: "Id", Width: utils.Min(50, idLength)},
			{Title: "ConnectionType", Width: utils.Min(50, utils.Max(6, len("ConnectionType")))},
		}

	} else {
		idLength := 6
		for _, r := range deviceList {
			idLength = utils.Max(idLength, len(r.SerialNumber))
			rows = append(rows, table.Row{r.SerialNumber, string(r.ConnectionType)})
		}
		for _, r := range deviceList {
			rows = append(rows, table.Row{r.SerialNumber, utils.FormatNetSpeed(r.ConnectionSpeed), string(r.ConnectionType)})
		}
		columns = []table.Column{
			{Title: "Id", Width: utils.Min(50, idLength)},
			{Title: "ConnectionSpeed", Width: utils.Min(50, utils.Max(6, len("ConnectionSpeed")))},
			{Title: "ConnectionType", Width: utils.Min(50, utils.Max(6, len("ConnectionType")))},
		}
	}
	tb := printer.NewTable(context.Background())
	if d.tableHeight > 0 {
		tb.SetHeight(d.tableHeight)
	}
	tb.SetColumns(columns).SetRows(rows)
	tb.Run()
	tb.Wait()
	return nil
}

func (d *DeviceHelper) HandleProxy(ctx context.Context, deviceId, protocol string, remote, local int) error {
	usbDriver := usbmuxd.NewUsbDriver()
	deviceList, err := usbDriver.DeviceList()
	if err != nil {
		log.Error(err)
	}
	if len(deviceList) == 0 {
		return errors.New("no device found")
	}
	var index int
	if deviceId == "" {
		index = 0
	} else {
		deviceFind := false
		for i, d := range deviceList {
			if strings.HasPrefix(d.SerialNumber, deviceId) {
				deviceFind = true
				index = i
				break
			}
		}
		if !deviceFind {
			return fmt.Errorf("device %s not found", deviceId)
		}
	}
	log.Debug("proxy info", "protocol", protocol, "remote", remote, "local", local, "deviceId", deviceList[index].SerialNumber, "index", index)
	proxyCtx, cancel := context.WithCancel(ctx)

	// Use errChan to receive possible errors.
	errChan := make(chan error, 1)
	go func() {
		err = usbmuxd.ProxyPort(proxyCtx, usbDriver, protocol, local, remote, index)
		if err != nil {
			log.Error(err)
			errChan <- err
		}
		close(errChan)
	}()

	// Listen for context cancellation signals.
	go func() {
		<-ctx.Done()
		cancel() // Ensure that when the parent context is canceled, the proxy context is also canceled.
	}()

	// Returns the first error that occurred
	select {
	case err := <-errChan:
		return err
	case <-time.After(time.Millisecond * 100): // 简单的延时确保代理已经启动
		return nil
	}
}
