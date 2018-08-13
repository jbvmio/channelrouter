package channelrouter

import (
	"reflect"
)

//ChanLink Here.
type ChanLink struct {
	Key      Key
	channel  chan interface{}
	refType  reflect.Type
	sent     uint32
	received uint32
}

//GetKey Here.
func (cl *ChanLink) GetKey() Key {
	return cl.Key
}

//Send Here.
func (cl *ChanLink) Send(i interface{}) {
	cl.channel <- i
}

//Receive Here.
func (cl *ChanLink) Receive() (interface{}, bool) {
	select {
	case s := <-cl.channel:
		return s, true
	default:
		return nil, false
	}
}

//SetType Here.
func (cl *ChanLink) SetType(t interface{}) {
	cl.refType = reflect.TypeOf(t)
}

//NewChanLink Here.
func NewChanLink(buffer int) *ChanLink {
	if buffer < 32 {
		buffer = 32
	}
	return &ChanLink{
		Key:     newKey(),
		channel: make(chan interface{}, buffer),
		refType: nil,
	}
}

func (cl *ChanLink) getAvailable() uint32 {
	return cl.sent - cl.received
}

func (cl *ChanLink) resetCounters() bool {
	if cl.getAvailable() == 0 {
		cl.sent = 0
		cl.received = 0
		return true
	}
	return false
}

func (cl *ChanLink) testTypeEq(t interface{}) bool {
	if reflect.TypeOf(t) == cl.refType {
		return true
	}
	return false
}

func (cl *ChanLink) typeNil(t interface{}) bool {
	if reflect.TypeOf(cl.refType) == nil {
		return true
	}
	return false
}

//Ideas ^
/*

Send / Receive Counters

*/
