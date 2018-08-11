package channelrouter

//ChanWriter Here.
type ChanWriter struct {
	key     Key
	cRouter *ChannelRouter
}

//Write Here.
func (cw *ChanWriter) Write(p []byte) (int, error) {
	var n int
	var err error
	for _, b := range p {
		cw.cRouter.Send(cw.key, b)
		n++
	}
	return n, err
}

//NewChanWriter returns a ChanWriter.
func (cr *ChannelRouter) NewChanWriter(k Key) *ChanWriter {
	return &ChanWriter{
		key:     k,
		cRouter: cr,
	}
}

//Send Here.
/*
func (cr *ChannelRouter) Send(k Key, i interface{}) {
	//cr.Channels[k].Send(i)
	cr.wg.Add(1)
	p := Packet{
		header: k,
		value:  i,
	}
	cr.ingress <- p
	cr.wg.Wait()
}

//NewChanWriter return an io.Writer for the corresponding channel referenced by k.
func (cr *ChannelRouter) NewChanWriter(k Key) ChanWriter {

	return ChanWriter{}
}


*/
