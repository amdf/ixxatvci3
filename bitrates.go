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
