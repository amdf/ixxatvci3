//go:build linux
// +build linux

package ixxatvci3

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

type ixxatConnection struct {
	name string // "can0"
	conn net.Conn
	tr   *socketcan.Transmitter
	recv *socketcan.Receiver
}

var devices map[uint8]*ixxatConnection

func init() {
	devices = make(map[uint8]*ixxatConnection)
}

//SelectDevice USB-to-CAN device select dialog.
//assignnumber - number to assign to the device.
// vcierr is 0 if there are no errors.
func SelectDevice(assignnumber uint8) (vcierr uint32) {
	return selectDevice(true, assignnumber)
}

//OpenDevice opens first USB-to-CAN device found.
//assignnumber - number to assign to the device.
// vcierr is 0 if there are no errors.
func OpenDevice(assignnumber uint8) (vcierr uint32) {
	return selectDevice(false, assignnumber)
}

func selectDevice(userselect bool, assignnumber uint8) (vcierr uint32) {
	if userselect {
		return VCI_E_NOT_IMPLEMENTED
	}
	ifaceName := fmt.Sprintf("can%d", assignnumber)
	devices[assignnumber] = &ixxatConnection{name: ifaceName}

	return
}

/*SetOperatingMode set operating mode at device with number "devnum".
Call it after SelectDevice but before OpenChannel.
11-bit mode is a default. Comma is a separator in opmode string.
opmode values:
	"11bit" or "standard" or "base",
	"29bit" or "extended",
	"err" or "errframe",
	"listen" or "listenonly" or "listonly",
	"low" or "lowspeed"
// vcierr is 0 if there are no errors.
*/
func SetOperatingMode(devnum uint8, opmode string) (vcierr uint32) {
	var setMode byte

	if strings.Contains(opmode, "11bit") || strings.Contains(opmode, "standard") || strings.Contains(opmode, "base") {
		setMode |= opmodeSTANDARD
	}
	if strings.Contains(opmode, "29bit") || strings.Contains(opmode, "extended") {
		setMode |= opmodeEXTENDED // reception of 29-bit id messages
	}
	if strings.Contains(opmode, "err") {
		setMode |= opmodeERRFRAME // enable reception of error frames
	}
	if strings.Contains(opmode, "list") {
		setMode |= opmodeLISTONLY // listen only mode (TX passive)
	}
	if strings.Contains(opmode, "low") {
		setMode |= opmodeLOWSPEED // use low speed bus interface
	}

	// HRESULT GOEXPORT CAN_VCI3_SetOperatingMode(UINT8 uDevNum, BYTE uCanOpMode)
	// ret := C.CAN_VCI3_SetOperatingMode(C.uchar(devnum), C.uchar(setMode))
	// vcierr = uint32(ret)

	return
}

//OpenChannel opens a channel on a previously opened device with devnum number, and btr0 and btr1 speed parameters.
//25 kbps is 0x1F 0x16.
//125 кб/с is 0x03 0x1C.
//vcierr is 0 if there are no errors.
func OpenChannel(devnum uint8, btr0 uint8, btr1 uint8) (vcierr uint32) {
	dev, ok := devices[devnum]
	if !ok {
		return VCI_E_NOT_INITIALIZED
	}

	var err error
	dev.conn, err = socketcan.DialContext(context.Background(), "can", dev.name)
	if err != nil {
		return VCI_E_FAIL
	}

	dev.tr = socketcan.NewTransmitter(dev.conn)
	dev.recv = socketcan.NewReceiver(dev.conn)

	return
}

// Send sends a data packet to device devnum.
// msgid - Identifier.
// rtr - Request flag, default value is false.
// msgdata - An array of 1 to 8 bytes. If rtr = true this field is ignored.
// vcierr is 0 if there are no errors.
func Send(devnum uint8, msgid uint32, rtr bool, msgdata []byte) (vcierr uint32) {

	dev, ok := devices[devnum]
	if !ok {
		return VCI_E_NOT_INITIALIZED
	}

	if len(msgdata) > 8 {
		return VCI_E_INVALIDARG
	}

	const maxmsgid11bit = 0x7FF
	const maxmsgid29bit = 0x1FFFFFFF

	ext := false
	if msgid > maxmsgid11bit {
		ext = true
	}

	fr := can.Frame{
		ID:         msgid & maxmsgid29bit,
		Length:     uint8(len(msgdata)),
		IsRemote:   rtr,
		IsExtended: ext,
	}

	copy(fr.Data[:], msgdata)

	err := dev.tr.TransmitFrame(context.Background(), fr)
	if err != nil {
		return VCI_E_FAIL
	}

	return
}

// Receive receives a message from a device with number "devnum".
// You need to call this function regularly so that the hardware message buffer does not overflow.
// Blocking call if no CAN messages are received.
// May return: VCI_E_OK, VCI_E_TIMEOUT, VCI_E_NO_DATA, VCI_E_INVALIDARG
func Receive(devnum uint8) (vcierr uint32, msgid uint32, rtr bool, msgdata [8]byte, msgdatasize uint8) {
	dev, ok := devices[devnum]
	if !ok {
		vcierr = VCI_E_NOT_INITIALIZED
		return
	}

	if !dev.recv.Receive() {
		vcierr = VCI_E_FAIL
		return
	}

	fr := dev.recv.Frame()
	msgid = fr.ID
	rtr = fr.IsRemote
	msgdatasize = fr.Length
	copy(msgdata[:], fr.Data[:])

	return
}

// GetStatus returns a structure containing various information about the connection status.
func GetStatus(devnum uint8) (status CANChanStatus, vcierr uint32) {
	return
}

//GetErrorText returns VCI error text by code
func GetErrorText(vcierr uint32) string {
	result := fmt.Sprint("error code", vcierr)

	return result
}

//CloseDevice close channel and free device with number "devnum".
func CloseDevice(devnum uint8) (vcierr uint32) {
	dev, ok := devices[devnum]
	if !ok {
		vcierr = VCI_E_NOT_INITIALIZED
		return
	}
	dev.conn.Close()
	delete(devices, devnum)
	return
}

// OpenChannelDetectBitrate opens a channel on the previously opened device with devnum number, and tries to determine the bitrate in the CAN channel.
// The bitrate is determined from the number of possible ones, specified through the bitrate array with several pairs of values for the btr0 and btr1 registers.
// If the channel is open and the bitrate is defined, a pair of values btr0, btr1 is returned.
func OpenChannelDetectBitrate(devnum uint8, timeout time.Duration, bitrate []BitrateRegisterPair) (detected BitrateRegisterPair, err error) {
	err = errors.New("not implemented")
	return
}

//see VCI canControlDetectBitrate
func openChannelDetectBitrate(devnum uint8, timeoutMs uint16, arrayBtr0 []byte, arrayBtr1 []byte) (vcierr uint32, indexArray int32) {
	vcierr = VCI_E_NOT_IMPLEMENTED
	return
}
