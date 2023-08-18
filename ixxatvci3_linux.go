//go:build linux
// +build linux

package ixxatvci3

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

type connectionCAN struct {
	name string // "can0"
	conn net.Conn
	tr   *socketcan.Transmitter
	recv *socketcan.Receiver
}

var devices map[uint8]*connectionCAN

func init() {
	devices = make(map[uint8]*connectionCAN)
}

// SelectDevice USB-to-CAN device select dialog.
// assignnumber - number to assign to the device.
// vcierr is 0 if there are no errors.
func SelectDevice(assignnumber uint8) (vcierr uint32) {
	return selectDevice(true, assignnumber)
}

// OpenDevice opens first USB-to-CAN device found.
// assignnumber - number to assign to the device.
// vcierr is 0 if there are no errors.
func OpenDevice(assignnumber uint8) (vcierr uint32) {
	return selectDevice(false, assignnumber)
}

func selectDevice(userselect bool, assignnumber uint8) (vcierr uint32) {
	if userselect {
		return VCI_E_NOT_IMPLEMENTED
	}

	_, ok := devices[assignnumber]
	if ok {
		return VCI_E_ALREADY_INITIALIZED
	}

	var dev connectionCAN
	dev.name = fmt.Sprintf("can%d", assignnumber)

	devices[assignnumber] = &dev

	return
}

/*
SetOperatingMode set operating mode at device with number "devnum".
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
	//NOTE: not implemented
	return
}

func execCmd(command string, args ...string) error {
	cmdObj := exec.Command(command, args...)
	cmdObj.Stdout = os.Stdout
	cmdObj.Stderr = os.Stderr
	return cmdObj.Run()
}

// OpenChannel opens a channel on a previously opened device with devnum number, and btr0 and btr1 speed parameters.
// 25 kbps is 0x1F 0x16.
// 125 кб/с is 0x03 0x1C.
// vcierr is 0 if there are no errors.
func OpenChannel(devnum uint8, btr0 uint8, btr1 uint8) (vcierr uint32) {
	dev, ok := devices[devnum]
	if !ok {
		return VCI_E_NOT_INITIALIZED
	}

	var err error
	speed := "0"
	brp := BitrateRegisterPair{Btr0: btr0, Btr1: btr1}
	switch brp {
	case Bitrate10kbps:
		speed = "10000"
	case Bitrate20kbps:
		speed = "20000"
	case Bitrate25kbps:
		speed = "25000"
	case Bitrate50kbps:
		speed = "50000"
	case Bitrate100kbps:
		speed = "100000"
	case Bitrate125kbps:
		speed = "125000"
	case Bitrate250kbps:
		speed = "250000"
	case Bitrate500kbps:
		speed = "500000"
	case Bitrate800kbps:
		speed = "800000"
	case Bitrate1000kbps:
		speed = "1000000"
	}

	err = execCmd("sudo", "ip", "link", "set", dev.name, "down")
	if err != nil {
		return VCI_E_FAIL
	}
	log.Println("link", dev.name, "restart")

	err = execCmd("sudo", "ip", "link", "set", dev.name, "up", "type", "can", "bitrate", speed)
	if err != nil {
		return VCI_E_FAIL
	}
	log.Println("link", dev.name, "up")

	dev.conn, err = socketcan.DialContext(context.Background(), "can", dev.name)
	if err != nil {
		return VCI_E_FAIL
	}

	dev.tr = socketcan.NewTransmitter(dev.conn)
	dev.recv = socketcan.NewReceiver(dev.conn)

	log.Println("connection", dev.name, "established")

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
	if nil == dev.tr {
		vcierr = VCI_E_NOT_INITIALIZED
		return
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

	if nil == dev.recv {
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

// GetErrorText returns VCI error text by code
func GetErrorText(vcierr uint32) string {
	result := fmt.Sprint("err: ", vcierr)

	return result
}

// CloseDevice close channel and free device with number "devnum".
func CloseDevice(devnum uint8) (vcierr uint32) {
	dev, ok := devices[devnum]
	if !ok {
		vcierr = VCI_E_NOT_INITIALIZED
		return
	}
	if nil != dev.recv && nil != dev.tr {

		dev.recv.Close()
		dev.tr.Close()
		dev.conn.Close()
		log.Println("connection closed")
		err := execCmd("ip", "link", "set", dev.name, "down")
		if err != nil {
			log.Println("close link error", err)
		} else {
			log.Println("link closed")
		}
	}
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
