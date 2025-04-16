package usbmuxd_device

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/solywsh/go-forensic/utils/iosutils/usbmuxd-device/usbmuxd"
	"github.com/solywsh/go-forensic/utils/logger"

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
	ConnectionType  usbmuxd.ConnectionType
}

type USBHub struct {
	//proto usbmuxd.Protocol
}

func NewUSBHub() *USBHub {
	return &USBHub{}
}

func (c *USBHub) DeviceList() (usbDevices []USBDevice, err error) {
	var proto *usbmuxd.Protocol
	if proto, err = usbmuxd.NewProtocol(usbmuxd.NewDefaultRequestFrame(usbmuxd.MessageTypeDeviceList), usbmuxd.PacketProtocolPlist, usbmuxd.PacketTypePlistPayload); err != nil {
		return nil, err
	}

	if err = proto.SendPacket(); err != nil {
		return nil, err
	}

	var respPacket *usbmuxd.ResponsePacket
	if respPacket, err = proto.RecvPacket(); err != nil {
		return nil, err
	}

	var devList usbmuxd.DeviceListResponseFrame
	if _, err = plist.Unmarshal(respPacket.Packet, &devList); err != nil {
		return nil, err
	}

	usbDevices = make([]USBDevice, 0, len(devList.DeviceList))
	for i := range devList.DeviceList {
		var dev USBDevice
		if devList.DeviceList[i].Properties.ConnectionType == usbmuxd.ConnectionTypeUSB {
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
	var proto *usbmuxd.Protocol
	if proto, err = usbmuxd.NewProtocol(
		usbmuxd.NewDefaultRequestFrame(usbmuxd.MessageTypeListen),
		usbmuxd.PacketProtocolPlist, usbmuxd.PacketTypePlistPayload); err != nil {
		return nil, err
	}
	if err = proto.SendPacket(); err != nil {
		return nil, err
	}

	C = make(chan USBDevice)
	go func() {
		defer close(C)
		for {
			var respPacket *usbmuxd.ResponsePacket
			if respPacket, err = proto.RecvPacket(); err != nil {
				break
			}

			var devInfo usbmuxd.DeviceResponseFrame
			if _, err = plist.Unmarshal(respPacket.Packet, &devInfo); err != nil {
				break
			}
			if devInfo.MessageType != string(usbmuxd.MessageTypeDeviceAdd) {
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
	var proto *usbmuxd.Protocol
	if proto, err = usbmuxd.NewProtocol(usbmuxd.NewConnectRequestFrame(devID, port), usbmuxd.PacketProtocolPlist, usbmuxd.PacketTypePlistPayload); err != nil {
		return nil, err
	}

	if err = proto.SendPacket(); err != nil {
		return nil, err
	}

	var respPacket *usbmuxd.ResponsePacket
	if respPacket, err = proto.RecvPacket(); err != nil {
		return nil, err
	}

	if respPacket.MsgType != usbmuxd.MessageTypeResult {
		return nil, fmt.Errorf("message type mismatch: expected '%s', got '%s'", usbmuxd.MessageTypeResult, respPacket.MsgType)
	}

	var result usbmuxd.ResultResponseFrame
	if _, err = plist.Unmarshal(respPacket.Packet, &result); err != nil {
		return nil, err
	}

	if result.ReplyCode != usbmuxd.ReplyCodeOK {
		return nil, fmt.Errorf("connect: %s", result.ReplyCode)
	}

	conn = proto.Conn()

	return
}

// ProxyPort Create a universal network proxy that forwards from a local port to a target port on an iOS device.
// Parameter description:
//   - usbHub: USB connection manager, if it is nil, a new one will be created.
//   - protocol: Network protocol type, "tcp" or "udp"
//   - localPort: Local listening port
//   - remotePort: Target port on iOS devices
//   - deviceIndex: Device index, specifies the device to be used, default is 0 (the first device).
func ProxyPort(usbHub *USBHub, protocol string, localPort int, remotePort int, deviceIndex ...int) error {
	if usbHub == nil {
		usbHub = NewUSBHub()
	}

	// Verify protocol type
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("unsupported protocol type: %s, only tcp or udp are supported", protocol)
	}

	// Get the list of connected devices.
	devices, err := usbHub.DeviceList()
	if err != nil {
		return fmt.Errorf("failed to obtain device list: %w", err)
	}

	if len(devices) == 0 {
		return fmt.Errorf("no connected devices")
	}

	// Determine the device index to be used.
	idx := 0
	if len(deviceIndex) > 0 && deviceIndex[0] >= 0 && deviceIndex[0] < len(devices) {
		idx = deviceIndex[0]
	}

	// Use the device with the specified index.
	device := devices[idx]
	log.Debug("using device: %s (ID: %d)", device.SerialNumber, device.DeviceID)

	// Start local network monitoring.
	listener, err := net.Listen(protocol, fmt.Sprintf(":%d", localPort))
	if err != nil {
		return fmt.Errorf("failed to start local %s listener: %w", protocol, err)
	}
	defer listener.Close()

	log.Debugf("agent has started: local %s port %d -> device %s port %d",
		protocol, localPort, device.SerialNumber, remotePort)

	for {
		// Accept local client connections
		localConn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept local connection: %w", err)
		}

		// Target port for connecting to the device
		remoteConn, err := usbHub.CreateConnect(device.DeviceID, remotePort)
		if err != nil {
			localConn.Close()
			return fmt.Errorf("failed to connect to device port %d: %w", remotePort, err)
		}

		log.Debugf("new connection: %s -> %s", localConn.RemoteAddr(), device.SerialNumber)

		// Start bidirectional data replication.
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

// ProxySSH proxies SSH from a local port to the iOS device's SSH port (22).
func ProxySSH(usbHub *USBHub, localPort int) error {
	// Directly call the general proxy function, specifying the SSH port (22).
	return ProxyPort(usbHub, "tcp", localPort, 22)
}
