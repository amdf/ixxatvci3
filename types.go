package ixxatvci3

import "C"

//BitrateRegisterPair Two CAN bitrate registers.
type BitrateRegisterPair struct {
	Btr0 byte
	Btr1 byte
}

//Bitrate predefinitions:
var (
	Bitrate10kbps   = BitrateRegisterPair{Btr0: 0x31, Btr1: 0x1C}
	Bitrate20kbps   = BitrateRegisterPair{Btr0: 0x18, Btr1: 0x1C}
	Bitrate25kbps   = BitrateRegisterPair{Btr0: 0x1F, Btr1: 0x16}
	Bitrate50kbps   = BitrateRegisterPair{Btr0: 0x09, Btr1: 0x1C}
	Bitrate100kbps  = BitrateRegisterPair{Btr0: 0x04, Btr1: 0x1C}
	Bitrate125kbps  = BitrateRegisterPair{Btr0: 0x03, Btr1: 0x1C}
	Bitrate250kbps  = BitrateRegisterPair{Btr0: 0x01, Btr1: 0x1C}
	Bitrate500kbps  = BitrateRegisterPair{Btr0: 0x00, Btr1: 0x1C}
	Bitrate800kbps  = BitrateRegisterPair{Btr0: 0x00, Btr1: 0x16}
	Bitrate1000kbps = BitrateRegisterPair{Btr0: 0x00, Btr1: 0x14}
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
