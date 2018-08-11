package channelrouter

import "sync"

//ChannelRouter Here.
type ChannelRouter struct {
	ingress  chan Packet
	run      bool
	running  bool
	Channels map[Key]*ChanLink
	wg       sync.WaitGroup
}

//AddChannel Here.
func (cr *ChannelRouter) AddChannel(buffer int) Key {
	cl := NewChanLink(buffer)
	cr.Channels[cl.Key] = cl
	return cl.Key
}

//Stop Here.
func (cr *ChannelRouter) Stop() {
	cr.run = false
}

//Send Here.
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

//Receive Here.
func (cr *ChannelRouter) Receive(k Key) Packet {
	r, ok := cr.Channels[k].Receive()
	p := Packet{
		header: k,
	}
	if !ok {
		p.value = nil
		return p
	}
	p.value = r
	return p
}

//Route Here.
func (cr *ChannelRouter) Route() {
	cr.run = true
	go func() {
		for {
			if cr.run {
				select {
				case p := <-cr.ingress:
					cr.Channels[p.header].Send(p.value)
					cr.wg.Done()
				default:
					cr.running = true
				}
				if !cr.run {
					return
				}
			}
		}
	}()
	for {
		if cr.running {
			break
		}
	}
}

//NewChannelRouter Here.
func NewChannelRouter() *ChannelRouter {
	channels := make(map[Key]*ChanLink)
	return &ChannelRouter{
		ingress:  make(chan Packet, 1024),
		Channels: channels,
	}
}
