package usbmuxd_device

import (
	"context"
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

type USBDriver struct {
	//proto usbmuxd.Protocol
}

func NewUsbDriver() *USBDriver {
	return &USBDriver{}
}

func (c *USBDriver) DeviceList() (usbDevices []USBDevice, err error) {
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

func (c *USBDriver) deviceListenAttached() (C chan USBDevice, err error) {
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

func (c *USBDriver) CreateConnect(devID int, port int) (conn net.Conn, err error) {
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
//   - ctx: Context for cancellation and timeout control.
//   - usbHub: USB connection manager, if it is nil, a new one will be created.
//   - protocol: Network protocol type, "tcp" or "udp"
//   - localPort: Local listening port
//   - remotePort: Target port on iOS devices
//   - deviceIndex: Device index, specifies the device to be used, default is 0 (the first device).
func ProxyPort(ctx context.Context, usbDriver *USBDriver, protocol string, localPort int, remotePort int, deviceIndex ...int) error {
	if usbDriver == nil {
		usbDriver = NewUsbDriver()
	}

	// Verify protocol type
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("unsupported protocol type: %s, only tcp or udp are supported", protocol)
	}

	// Get the list of connected devices.
	devices, err := usbDriver.DeviceList()
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
	config := &net.ListenConfig{}
	listener, err := config.Listen(ctx, protocol, fmt.Sprintf(":%d", localPort))
	if err != nil {
		return fmt.Errorf("failed to start local %s listener: %w", protocol, err)
	}

	// Ensure that the listener is closed when the function exits or the context is canceled.
	defer listener.Close()

	log.Debugf("agent has started: local %s port %d -> device %s port %d",
		protocol, localPort, device.SerialNumber, remotePort)

	// Create channels for handling connections and error channels.
	connChan := make(chan net.Conn)
	errChan := make(chan error, 1)

	// Accepting connections in the background.
	go func() {
		for {
			localConn, err := listener.Accept()
			if err != nil {
				// Check if there is an error due to the listener being closed.
				if errors.Is(err, net.ErrClosed) {
					return
				}
				errChan <- fmt.Errorf("failed to accept local connection: %w", err)
				return
			}
			connChan <- localConn
		}
	}()

	// The main loop handles connections and context cancelation
	for {
		select {
		case <-ctx.Done():
			// Context canceled, exiting gracefully.
			log.Debugf("context cancelled, stopping proxy from local port %d to device %s port %d",
				localPort, device.SerialNumber, remotePort)
			return ctx.Err()

		case err := <-errChan:
			// Error processing listener
			return err

		case localConn := <-connChan:
			// Successfully accepted the connection, creating forwarding.
			remoteConn, err := usbDriver.CreateConnect(device.DeviceID, remotePort)
			if err != nil {
				localConn.Close()
				log.Errorf("failed to connect to device port %d: %v", remotePort, err)
				continue
			}

			log.Debugf("new connection: %s -> %s", localConn.RemoteAddr(), device.SerialNumber)

			// Create a sub-context for each connection so that they can be canceled individually.
			connCtx, cancel := context.WithCancel(ctx)

			// Bidirectional replication data
			go func() {
				defer localConn.Close()
				defer remoteConn.Close()
				defer cancel() // Ensure to cancel the connection context.

				// Using the separate copy function, disconnections can be detected.
				buf := make([]byte, 32*1024)
				_, err := io.CopyBuffer(localConn, remoteConn, buf)
				if err != nil && !errors.Is(err, io.EOF) {
					log.Debugf("error copying data from device to local: %v", err)
				}
			}()

			go func() {
				defer localConn.Close()
				defer remoteConn.Close()
				defer cancel() // Ensure to cancel the connection context.

				// Using the separate copy function, the disconnection can be detected.
				buf := make([]byte, 32*1024)
				_, err := io.CopyBuffer(remoteConn, localConn, buf)
				if err != nil && !errors.Is(err, io.EOF) {
					log.Debugf("error copying data from local to device: %v", err)
				}
			}()

			// Listening for connection context cancellation
			go func() {
				<-connCtx.Done()
				// When the connection context is canceled, we don't need to do anything,
				// because the defer statement will close the connection.
			}()
		}
	}
}

// ProxySSH proxies SSH from a local port to the iOS device's SSH port (22).
func ProxySSH(usbHub *USBDriver, localPort int) error {
	// Directly call the general proxy function, specifying the SSH port (22).
	return ProxyPort(context.Background(), usbHub, "tcp", localPort, 22)
}
