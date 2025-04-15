package go_usbmuxd_device

import (
	"errors"
	"fmt"
	"github.com/solywsh/go-forensic/utils/iosutils/usbmux"
	"github.com/solywsh/go-forensic/utils/logger"
	"io"
	"net"

	"howett.net/plist"
)

var (
	log = logger.NewLogger()
)

var ErrNoFindUSBDevice = errors.New("no find match device")

type USBDevice struct {
	DeviceID        int
	LocationID      int
	ProductID       int
	SerialNumber    string
	ConnectionSpeed int
	ConnectionType  usbmux.ConnectionType
}

type USBHub struct {
	// proto usbmux.Protocol
}

func NewUSBHub() *USBHub {
	return &USBHub{}
}

func (c *USBHub) DeviceList() (usbDevices []USBDevice, err error) {
	var proto *usbmux.Protocol
	if proto, err = usbmux.NewProtocol(usbmux.NewDefaultRequestFrame(usbmux.MessageTypeDeviceList), usbmux.PacketProtocolPlist, usbmux.PacketTypePlistPayload); err != nil {
		return nil, err
	}

	if err = proto.SendPacket(); err != nil {
		return nil, err
	}

	var respPacket *usbmux.ResponsePacket
	if respPacket, err = proto.RecvPacket(); err != nil {
		return nil, err
	}

	var devList usbmux.DeviceListResponseFrame
	if _, err = plist.Unmarshal(respPacket.Packet, &devList); err != nil {
		return nil, err
	}

	usbDevices = make([]USBDevice, 0, len(devList.DeviceList))
	for i := range devList.DeviceList {
		var dev USBDevice
		if devList.DeviceList[i].Properties.ConnectionType == usbmux.ConnectionTypeUSB {
			dev.DeviceID = devList.DeviceList[i].Properties.DeviceID
			dev.LocationID = devList.DeviceList[i].Properties.LocationID
			dev.ProductID = devList.DeviceList[i].Properties.ProductID
			dev.SerialNumber = devList.DeviceList[i].Properties.SerialNumber
			dev.ConnectionSpeed = devList.DeviceList[i].Properties.ConnectionSpeed
			dev.ConnectionType = devList.DeviceList[i].Properties.ConnectionType
			usbDevices = append(usbDevices, dev)
		}
	}

	if len(usbDevices) == 0 {
		return nil, ErrNoFindUSBDevice
	}

	return
}

func (c *USBHub) deviceListenAttached() (C chan USBDevice, err error) {
	var proto *usbmux.Protocol
	if proto, err = usbmux.NewProtocol(
		usbmux.NewDefaultRequestFrame(usbmux.MessageTypeListen),
		usbmux.PacketProtocolPlist, usbmux.PacketTypePlistPayload); err != nil {
		return nil, err
	}
	if err = proto.SendPacket(); err != nil {
		return nil, err
	}

	C = make(chan USBDevice)
	go func() {
		defer close(C)
		for {
			var respPacket *usbmux.ResponsePacket
			if respPacket, err = proto.RecvPacket(); err != nil {
				break
			}

			var devInfo usbmux.DeviceResponseFrame
			if _, err = plist.Unmarshal(respPacket.Packet, &devInfo); err != nil {
				break
			}
			if devInfo.MessageType != string(usbmux.MessageTypeDeviceAdd) {
				continue
			}
			var dev USBDevice
			dev.DeviceID = devInfo.Properties.DeviceID
			dev.LocationID = devInfo.Properties.LocationID
			dev.ProductID = devInfo.Properties.ProductID
			dev.SerialNumber = devInfo.Properties.SerialNumber
			dev.ConnectionSpeed = devInfo.Properties.ConnectionSpeed
			dev.ConnectionType = devInfo.Properties.ConnectionType
			C <- dev
		}
	}()
	return C, nil
}

func (c *USBHub) CreateConnect(devID int, port int) (conn net.Conn, err error) {
	var proto *usbmux.Protocol
	if proto, err = usbmux.NewProtocol(usbmux.NewConnectRequestFrame(devID, port), usbmux.PacketProtocolPlist, usbmux.PacketTypePlistPayload); err != nil {
		return nil, err
	}

	if err = proto.SendPacket(); err != nil {
		return nil, err
	}

	var respPacket *usbmux.ResponsePacket
	if respPacket, err = proto.RecvPacket(); err != nil {
		return nil, err
	}

	if respPacket.MsgType != usbmux.MessageTypeResult {
		return nil, fmt.Errorf("message type mismatch: expected '%s', got '%s'", usbmux.MessageTypeResult, respPacket.MsgType)
	}

	var result usbmux.ResultResponseFrame
	if _, err = plist.Unmarshal(respPacket.Packet, &result); err != nil {
		return nil, err
	}

	if result.ReplyCode != usbmux.ReplyCodeOK {
		return nil, fmt.Errorf("connect: %s", result.ReplyCode)
	}

	conn = proto.Conn()

	return
}

func Debug(b ...bool) {
	if len(b) == 0 {
		b = []bool{true}
	}
	usbmux.Debug = b[0]
}

// ProxySSH proxies SSH from a local port to the iOS device's SSH port (22).
func ProxySSH(usbHub *USBHub, localPort int) error {
	if usbHub == nil {
		usbHub = NewUSBHub()
	}

	// Get the list of connected devices
	devices, err := usbHub.DeviceList()
	if err != nil {
		return fmt.Errorf("failed to get device list: %w", err)
	}

	if len(devices) == 0 {
		return fmt.Errorf("no devices connected")
	}

	// Use the first device for the proxy
	device := devices[0]
	log.Print(device)
	// Start a local TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		return fmt.Errorf("failed to start local listener: %w", err)
	}
	defer listener.Close()

	fmt.Printf("Proxying SSH from local port %d to device %s:22\n", localPort, device.SerialNumber)

	for {
		// Accept a connection from the local client
		localConn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept local connection: %w", err)
		}

		// Connect to the device's SSH port (22)
		remoteConn, err := usbHub.CreateConnect(device.DeviceID, 22)
		if err != nil {
			localConn.Close()
			return fmt.Errorf("failed to connect to device SSH port: %w", err)
		}

		// Start bidirectional copy between local and remote connections
		go func() {
			defer localConn.Close()
			defer remoteConn.Close()
			io.Copy(localConn, remoteConn)
		}()

		go func() {
			defer localConn.Close()
			defer remoteConn.Close()
			io.Copy(remoteConn, localConn)
		}()
	}
}
