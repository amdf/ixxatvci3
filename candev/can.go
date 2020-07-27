package candev

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/amdf/ixxatvci3"
)

//Device is a USB-to-CAN device type
type Device struct {
	number                 uint8
	stop                   bool
	canMessagesChannel     chan Message
	canAdditionalChannels  map[uint]chan Message
	iChIndex               uint
	muAddCh                sync.Mutex
	muReadCAN              sync.Mutex
	bReadCANBackground     bool
	RcvOkCount             uint
	RcvErrCount            uint
	RcvProcessedBackground uint
	RcvProcessedActive     uint
	RcvBackgroundNoData    uint
}

//Message is a CAN message
type Message struct {
	ID   uint32
	Rtr  bool
	Ext  bool //true if 29-bit mode
	Len  uint8
	Data [8]byte
}

func (dev *Device) canReaderThread() {
	for !dev.stop {

		var vcierr uint32
		var rxMsg Message

		vcierr, rxMsg.ID, rxMsg.Rtr, rxMsg.Data, rxMsg.Len =
			ixxatvci3.Receive(dev.number)

		if 0 == vcierr {
			dev.RcvOkCount++

			dev.canMessagesChannel <- rxMsg //blocking call

			//sending to additional channels
			dev.muAddCh.Lock()
			for _, addch := range dev.canAdditionalChannels {
				addch <- rxMsg
			}
			dev.muAddCh.Unlock()
		} else {
			dev.RcvErrCount++
		}
		runtime.Gosched()
	}
}

//nonblocking background reading
func (dev *Device) canBackgroundReaderThread() {
	for !dev.stop {
		if dev.bReadCANBackground {
			select {
			case <-dev.canMessagesChannel:
				dev.RcvProcessedBackground++

			default:
				dev.RcvBackgroundNoData++

			}
		}
		runtime.Gosched()
	}
}

//GetMsgByID get msg from CAN with id
func (dev *Device) GetMsgByID(id uint32, timeout time.Duration) (msg Message, err error) {
	msg, err = dev.GetMsgByIDList(map[uint32]bool{id: true}, timeout)
	return
}

//GetMsgByIDList waits for msg from CAN with id list
func (dev *Device) GetMsgByIDList(idlist map[uint32]bool, timeout time.Duration) (msg Message, err error) {
	if nil == dev {
		err = fmt.Errorf("%s", "Device == nil")
		return
	}
	dev.muReadCAN.Lock()
	defer dev.muReadCAN.Unlock()
	msgReceived := 0
	dev.disableBackgroundCAN()      //disable background read
	defer dev.enableBackgroundCAN() //re-enable afterwards
	start := time.Now()
	for time.Since(start) <= timeout {
		select {
		case rmsg := <-dev.canMessagesChannel:
			dev.RcvProcessedActive++
			msgReceived++
			_, ok := idlist[rmsg.ID]
			if ok {
				msg = rmsg
				return
			}
		default:

		}
		runtime.Gosched()
	}

	err = fmt.Errorf("timeout (%d msgs)", msgReceived)
	return
}

func (dev *Device) disableBackgroundCAN() {
	dev.bReadCANBackground = false
}

func (dev *Device) enableBackgroundCAN() {
	dev.bReadCANBackground = true
}

//GetMsgByIDAndSize waits for msg from CAN with id and size
func (dev *Device) GetMsgByIDAndSize(id uint32, size uint8, timeout time.Duration) (msg Message, err error) {
	if nil == dev {
		err = fmt.Errorf("%s", "Device == nil")
		return
	}
	dev.muReadCAN.Lock()
	defer dev.muReadCAN.Unlock()
	msgReceived := 0
	dev.disableBackgroundCAN()      //disable background read
	defer dev.enableBackgroundCAN() //re-enable afterwards
	start := time.Now()
	for time.Since(start) <= timeout {
		select {
		case rmsg := <-dev.canMessagesChannel:
			dev.RcvProcessedActive++
			msgReceived++
			if rmsg.ID == id && rmsg.Len == size {
				msg = rmsg
				return
			}
		default:

		}
		runtime.Gosched()
	}

	err = fmt.Errorf("timeout (%d msgs)", msgReceived)
	return
}

//GetMsgRTR waits for msg from CAN with id and RTR flag set
func (dev *Device) GetMsgRTR(id uint32, timeout time.Duration) (ok bool, err error) {
	if nil == dev {
		err = fmt.Errorf("%s", "Device == nil")
		return
	}
	dev.muReadCAN.Lock()
	defer dev.muReadCAN.Unlock()
	msgReceived := 0
	dev.disableBackgroundCAN()      //disable background read
	defer dev.enableBackgroundCAN() //re-enable afterwards
	start := time.Now()
	for time.Since(start) <= timeout {
		select {
		case rmsg := <-dev.canMessagesChannel:
			dev.RcvProcessedActive++
			msgReceived++
			if rmsg.ID == id && rmsg.Rtr == true {
				ok = true
				return
			}
		default:

		}
		runtime.Gosched()
	}

	err = fmt.Errorf("timeout (%d msgs)", msgReceived)
	return
}

//Run starts receiving
func (dev *Device) Run() {
	go dev.canReaderThread()
	go dev.canBackgroundReaderThread()
}

func (dev *Device) deviceInit(devNum uint8) {
	dev.number = devNum
	dev.bReadCANBackground = true
	dev.canMessagesChannel = make(chan Message)
	dev.stop = false
	dev.canAdditionalChannels = make(map[uint]chan Message)
}

//Init first USB-to-CAN device found
//btr0, btr1 - CAN speed register values.
func (dev *Device) Init(btr0, btr1 uint8) (err error) {
	if nil == dev {
		err = fmt.Errorf("%s", "null ptr")
		return
	}

	var b Builder
	_, err = b.Mode("11bit").Speed(ixxatvci3.BitrateRegisterPair{Btr0: btr0, Btr1: btr1}).Get()
	dev.deviceInit(0)

	return
}

//InitSelect shows device selection dialog and initializes selected deviсe.
//devNum - device number to assign
//btr0, btr1 - CAN speed register values.
func (dev *Device) InitSelect(devNum, btr0, btr1 uint8) (err error) {
	if nil == dev {
		err = fmt.Errorf("%s", "null ptr")
		return
	}

	var b Builder
	_, err = b.Number(devNum).Mode("11bit").Speed(ixxatvci3.BitrateRegisterPair{Btr0: btr0, Btr1: btr1}).SelectDevice(true).Get()
	dev.deviceInit(devNum)

	return
}

//InitSelectDetectBitrate steps:
//Shows device selection dialog.
//Initializes device and assign devNum.
//Detects bitrate from list of possible bitrates.
//if success, returns detected bitrate.
func (dev *Device) InitSelectDetectBitrate(devNum uint8, timeout time.Duration, bitrate []ixxatvci3.BitrateRegisterPair) (detected ixxatvci3.BitrateRegisterPair, err error) {

	if nil == dev {
		err = fmt.Errorf("%s", "null ptr")
		return
	}

	var b Builder
	_, err = b.Number(devNum).Mode("11bit").Timeout(timeout).Detect(bitrate).SelectDevice(true).Get()
	dev.deviceInit(devNum)

	return
}

//Stop stops receiving
func (dev *Device) Stop() {
	if nil == dev {
		return
	}
	dev.stop = true
	time.Sleep(2 * time.Second)
	for _, addch := range dev.canAdditionalChannels {
		close(addch)
	}
	ixxatvci3.CloseDevice(dev.number)
}

//GetBusLoad CAN bus load %
func (dev *Device) GetBusLoad(nDev uint8) uint8 {
	if nil == dev {
		return 0
	}
	st, _ := ixxatvci3.GetStatus(nDev)
	return st.LineStatus.BusLoad
}

/*Send msg to CAN.
If msg.ID <= 0x7FF, msg is 11-bit.
If msg.ID > 0x7FF, msg is 29-bit.

For sending small IDs in 29 bit mode use SendExt.
*/
func (dev *Device) Send(msg Message) (err error) {
	if nil == dev {
		err = fmt.Errorf("%s", "null ptr")
		return
	}
	var vcierr uint32
	if !msg.Ext {
		vcierr = ixxatvci3.Send(dev.number, msg.ID, msg.Rtr, msg.Data[0:msg.Len])
	} else {
		//hi bit is 29-bit mode flag for ixxatvci3 package
		vcierr = ixxatvci3.Send(dev.number, msg.ID|(1<<31), msg.Rtr, msg.Data[0:msg.Len])
	}
	if ixxatvci3.VCI_OK != vcierr {
		err = fmt.Errorf("%s", ixxatvci3.GetErrorText(vcierr))
	}
	return
}

//GetMsgChannelCopy returns channel with all messages received from CAN.
//idx - channel index for use with CloseMsgChannelCopy().
func (dev *Device) GetMsgChannelCopy() (ch <-chan Message, idx uint) {
	if nil == dev {
		return
	}
	dev.canAdditionalChannels[dev.iChIndex] = make(chan Message)
	ch = dev.canAdditionalChannels[dev.iChIndex]
	idx = dev.iChIndex
	dev.iChIndex++

	return
}

//CloseMsgChannelCopy close msg channel.
//idx is a channel index from GetMsgChannelCopy().
func (dev *Device) CloseMsgChannelCopy(idx uint) {
	if nil == dev {
		return
	}

	_, ok := dev.canAdditionalChannels[idx]
	if ok {
		dev.muAddCh.Lock()
		close(dev.canAdditionalChannels[idx])
		delete(dev.canAdditionalChannels, idx)
		dev.muAddCh.Unlock()
	}
}
