package candev

import (
	"fmt"
	"time"

	"github.com/amdf/ixxatvci3"
)

type builder interface {
	Get() (dev *Device, err error)
}

//Builder for configuring candev.Device
type Builder struct {
	builder
	dev             Device
	speed           ixxatvci3.BitrateRegisterPair
	wantBitrateList []ixxatvci3.BitrateRegisterPair
	foundBitrate    ixxatvci3.BitrateRegisterPair
	detectTimeout   time.Duration
	mode            string
	selectDevice    bool
	detectBitrate   bool
	number          uint8
}

//Get candev.Device
func (b *Builder) Get() (dev *Device, err error) {
	var vcierr uint32

	defer func() {
		if ixxatvci3.VCI_OK != vcierr {
			err = fmt.Errorf("candev.Builder:%s", ixxatvci3.GetErrorText(vcierr))
		}

		if nil == err {
			dev = &b.dev
			dev.deviceInit(b.number)
		}
	}()

	if b.selectDevice {
		vcierr = ixxatvci3.SelectDevice(b.number)
	} else {
		vcierr = ixxatvci3.OpenDevice(b.number)
	}
	if ixxatvci3.VCI_OK != vcierr {
		return
	}
	if b.mode == "" {
		b.mode = "11bit"
	}
	vcierr = ixxatvci3.SetOperatingMode(b.number, b.mode)
	if ixxatvci3.VCI_OK != vcierr {
		return
	}
	if b.detectBitrate {
		if 0 == b.detectTimeout {
			b.detectTimeout = 5 * time.Second
		}
		b.foundBitrate, err = ixxatvci3.OpenChannelDetectBitrate(b.number, b.detectTimeout, b.wantBitrateList)
		if err != nil {
			err = fmt.Errorf("candev.Builder:bitrate detect failed:%s", err.Error())
			return
		}
		foundInList := false
		for _, x := range b.wantBitrateList {
			if x == b.foundBitrate {
				foundInList = true
				break
			}
		}
		if !foundInList {
			err = fmt.Errorf("candev.Builder:cannot find desired bitrate")
			return
		}
	} else {
		vcierr = ixxatvci3.OpenChannel(b.number, b.speed.Btr0, b.speed.Btr1)
	}

	return
}

//Timeout for bitrate detection. Only use with Detect or AutoDetect.
func (b *Builder) Timeout(t time.Duration) *Builder {
	b.detectTimeout = t
	return b
}

//Number set device number (for multi-device configuration).
func (b *Builder) Number(number uint8) *Builder {
	b.number = number
	return b
}

//Speed set speed.
func (b *Builder) Speed(pair ixxatvci3.BitrateRegisterPair) *Builder {
	b.speed = pair
	return b
}

//Btr0 register BTR0. Use with Btr1()
func (b *Builder) Btr0(val byte) *Builder {
	b.speed.Btr0 = val
	return b
}

//Btr1 register BTR1. Use with Btr0()
func (b *Builder) Btr1(val byte) *Builder {
	b.speed.Btr1 = val
	return b
}

/*Mode set device mode.
Possible values:
"11bit" or "standard" or "base",
"29bit" or "extended",
"err" or "errframe",
"listen" or "listenonly" or "listonly",
"low" or "lowspeed".
Default is "11bit".
*/
func (b *Builder) Mode(mode string) *Builder {
	b.mode = mode
	return b
}

//SelectDevice shows dialog if true
func (b *Builder) SelectDevice(selectDevice bool) *Builder {
	b.selectDevice = selectDevice
	return b
}

//Detect bitrate from list of possble bitrates
func (b *Builder) Detect(bitrates []ixxatvci3.BitrateRegisterPair) *Builder {
	b.wantBitrateList = bitrates
	b.detectBitrate = true
	return b
}

//AutoDetect from common bitrates
func (b *Builder) AutoDetect() *Builder {
	b.wantBitrateList = []ixxatvci3.BitrateRegisterPair{
		ixxatvci3.Bitrate10kbps,
		ixxatvci3.Bitrate20kbps,
		ixxatvci3.Bitrate25kbps,
		ixxatvci3.Bitrate50kbps,
		ixxatvci3.Bitrate100kbps,
		ixxatvci3.Bitrate125kbps,
		ixxatvci3.Bitrate250kbps,
		ixxatvci3.Bitrate500kbps,
		ixxatvci3.Bitrate800kbps,
		ixxatvci3.Bitrate1000kbps,
	}
	b.detectBitrate = true
	return b
}
