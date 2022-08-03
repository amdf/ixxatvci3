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

const vciMaxErrStrLen = 256 // maximum length of an error string

//CANLineStatus информация о статусе линии CAN
type CANLineStatus struct {
	OpMode  uint8  // current CAN operating mode
	BtReg0  uint8  // current bus timing register 0 value
	BtReg1  uint8  // current bus timing register 1 value
	BusLoad uint8  // average bus load in percent (0..100)
	Status  uint32 // status of the CAN controller (see CAN_STATUS_)
}

//CANChanStatus информация о статусе канала связи
type CANChanStatus struct {
	LineStatus CANLineStatus // current CAN line status
	Activated  uint32        // TRUE if the channel is activated
	RxOverrun  uint32        // TRUE if receive FIFO overrun occurs
	RxFifoLoad uint8         // receive FIFO load in percent (0..100)
	TxFifoLoad uint8         // transmit FIFO load in percent (0..100)
}

const (
	opmodeUNDEFINED = 0x00 // undefined
	opmodeSTANDARD  = 0x01 // reception of 11-bit id messages
	opmodeEXTENDED  = 0x02 // reception of 29-bit id messages
	opmodeERRFRAME  = 0x04 // enable reception of error frames
	opmodeLISTONLY  = 0x08 // listen only mode (TX passive)
	opmodeLOWSPEED  = 0x10 // use low speed bus interface
)

type ixxatConnection struct {
	name string // "can0"
	conn net.Conn
	tr   *socketcan.Transmitter
	recv *socketcan.Receiver
}

var devices map[uint8]ixxatConnection

func init() {
	devices = make(map[uint8]ixxatConnection)
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
	devices[assignnumber] = ixxatConnection{name: ifaceName}
	// var us uint8

	// if userselect {
	// 	us = 1
	// }
	// HRESULT CAN_VCI3_SelectDevice(UINT8 bUserSelect, UINT8 uAssignNumber);
	// ret := C.CAN_VCI3_SelectDevice(
	// 	C.uchar(us),
	// 	C.uchar(assignnumber))
	// vcierr = uint32(ret)
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

	// HRESULT CAN_VCI3_OpenConnection(UINT8 uDevNum, UINT8 uBtr0, UINT8 uBtr1);
	// ret := C.CAN_VCI3_OpenConnection(
	// 	C.uchar(devnum),
	// 	C.uchar(btr0),
	// 	C.uchar(btr1))
	// vcierr = uint32(ret)
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

	// var irtr uint8

	// // HRESULT CAN_VCI3_RxData(UINT8 uDevNum, UINT32 * uMsgId, UINT8 * bRtr, BYTE * MsgData, UINT8 * uMsgDataSize);
	// ret := C.CAN_VCI3_RxData(
	// 	C.uchar(devnum),
	// 	(*C.uint)(unsafe.Pointer(&msgid)),
	// 	(*C.uchar)(unsafe.Pointer(&irtr)),
	// 	(*C.uchar)(unsafe.Pointer(&msgdata[0])),
	// 	(*C.uchar)(unsafe.Pointer(&msgdatasize)))
	// vcierr = uint32(ret)

	// rtr = bool(irtr != 0)

	return
}

// GetStatus returns a structure containing various information about the connection status.
func GetStatus(devnum uint8) (status CANChanStatus, vcierr uint32) {

	// HRESULT CAN_VCI3_GetStatus(UINT8 uDevNum, PCANCHANSTATUS pCanStat)
	// ret := C.CAN_VCI3_GetStatus(C.uchar(devnum), (*C.CANCHANSTATUS)(unsafe.Pointer(&status)))
	// vcierr = uint32(ret)

	return
}

//GetErrorText returns VCI error text by code
func GetErrorText(vcierr uint32) string {
	//void CAN_VCI3_FormatError(HRESULT hrError, PCHAR pszText, UINT32 dwSize)
	// buf := make([]C.char, vciMaxErrStrLen)

	// C.vciFormatError(C.long(vcierr), &buf[0], C.uint(vciMaxErrStrLen))

	// result := C.GoString(&buf[0])
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
	// HRESULT CAN_VCI3_CloseDevice(UINT8 uDevNum);
	// ret := C.CAN_VCI3_CloseDevice(C.uchar(devnum))
	// vcierr = uint32(ret)
	return
}

// OpenChannelDetectBitrate opens a channel on the previously opened device with devnum number, and tries to determine the bitrate in the CAN channel.
// The bitrate is determined from the number of possible ones, specified through the bitrate array with several pairs of values for the btr0 and btr1 registers.
// If the channel is open and the bitrate is defined, a pair of values btr0, btr1 is returned.
func OpenChannelDetectBitrate(devnum uint8, timeout time.Duration, bitrate []BitrateRegisterPair) (detected BitrateRegisterPair, err error) {
	err = errors.New("not implemented")
	// if len(bitrate) <= 0 {
	// 	err = fmt.Errorf("%s", "bitrate array is empty")
	// 	return
	// }

	// var buf0, buf1 bytes.Buffer

	// for _, b := range bitrate {
	// 	buf0.WriteByte(b.Btr0)
	// 	buf1.WriteByte(b.Btr1)
	// }

	// timeoutMs := uint16(timeout / time.Millisecond)

	// vcierr, indexArray := openChannelDetectBitrate(devnum, timeoutMs, buf0.Bytes(), buf1.Bytes())
	// if C.VCI_OK != vcierr {
	// 	err = fmt.Errorf("%s", GetErrorText(vcierr))
	// 	return
	// }

	// if (indexArray < 0) || int(indexArray) >= len(bitrate) {
	// 	err = fmt.Errorf("%s", "wrong index of bitrate array")
	// 	return
	// }

	// detected = bitrate[indexArray]

	return
}

//see VCI canControlDetectBitrate
func openChannelDetectBitrate(devnum uint8, timeoutMs uint16, arrayBtr0 []byte, arrayBtr1 []byte) (vcierr uint32, indexArray int32) {
	vcierr = VCI_E_NOT_IMPLEMENTED
	// len1 := len(arrayBtr0)
	// len2 := len(arrayBtr1)
	// if len1 != len2 {
	// 	vcierr = C.VCI_E_INVALIDARG
	// 	return
	// }

	// arrayElementCount := uint32(len1)

	// // HRESULT CAN_VCI3_OpenConnectionDetectBitrate(UINT8 uDevNum, UINT16 uTimeoutMs, UINT32 uArrayElementCount, BYTE * ArrayBtr0, BYTE * ArrayBtr1, INT32 * pIndexArray) {
	// ret := C.CAN_VCI3_OpenConnectionDetectBitrate(
	// 	C.uchar(devnum),
	// 	C.ushort(timeoutMs),
	// 	C.uint(arrayElementCount),
	// 	(*C.uchar)(unsafe.Pointer(&arrayBtr0[0])),
	// 	(*C.uchar)(unsafe.Pointer(&arrayBtr1[0])),
	// 	(*C.int)(unsafe.Pointer(&indexArray)))

	// vcierr = uint32(ret)

	return
}
