package channelrouter

//ChanLink Here.
type ChanLink struct {
	Key     Key
	channel chan interface{}
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

//NewChanLink Here.
func NewChanLink(buffer int) *ChanLink {
	if buffer < 32 {
		buffer = 32
	}
	return &ChanLink{
		Key:     newKey(),
		channel: make(chan interface{}, buffer),
	}
}
