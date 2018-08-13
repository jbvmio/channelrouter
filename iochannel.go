package channelrouter

import "fmt"

//MakeIoChannel returns a IoChannel for the corresponding channel referenced by k.
func (cr *ChannelRouter) MakeIoChannel(k Key) *IoChannel {
	cr.SetType(k, byte(0))
	return &IoChannel{
		key:     k,
		cRouter: cr,
	}
}

//IoChannel Here.
type IoChannel struct {
	key     Key
	cRouter *ChannelRouter
}

//Write Here.
func (cw *IoChannel) Write(p []byte) (int, error) {
	var n int
	var err error
	for _, b := range p {
		cw.cRouter.Send(cw.key, b)
		n++
	}
	if n != len(p) {
		err = fmt.Errorf("EOF: expecting %v bytes, received %v", len(p), n)
	}
	return n, err
}

/*
//Read Here.
func (cw *IoChannel) Read(p []byte) (int, error) {
	var n int
	var err error
	for range p {
		cw.cRouter.Receive(cw.key)
		n++
	}
	if n == 0 {
		err = fmt.Errorf("Received EOF")
	}
	return n, err
}
*/
