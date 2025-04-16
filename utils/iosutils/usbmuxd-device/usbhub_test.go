package usbmuxd_device

import (
	"fmt"
	"github.com/solywsh/go-forensic/utils/logger"
	"net"
	"testing"
	"time"
)

func TestUSBHub_DeviceList(t *testing.T) {
	usbHub := NewUSBHub()

	_ = logger.NewDebugLogger()

	devices, err := usbHub.DeviceList()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(devices)

	conn, err := usbHub.CreateConnect(devices[0].DeviceID, 8100)
	if err != nil {
		t.Fatal(err)
	}
	_ = conn
}

func TestDeviceList(t *testing.T) {
	u := NewUSBHub()
	deviceList, err := u.DeviceList()
	if err != nil {
		t.Error(err)
		return
	}
	for _, device := range deviceList {
		t.Log(device.SerialNumber)
	}
}

func TestProxySSH(t *testing.T) {
	usbHub := NewUSBHub()

	// Use a random local port for testing
	localPort := 2222

	// Run the ProxySSH function in a separate goroutine
	go func() {
		err := ProxySSH(usbHub, localPort)
		if err != nil {
			t.Errorf("ProxySSH failed: %v", err)
		}
	}()

	// Give the proxy some time to start
	time.Sleep(3 * time.Second)

	// Try to connect to the local port
	conn, err := net.Dial("tcp", "127.0.0.1:2222")
	if err != nil {
		t.Fatalf("Failed to connect to local proxy: %v", err)
	}
	defer conn.Close()

	// If connection is successful, the test passes
	t.Log("Successfully connected to local proxy")
}
