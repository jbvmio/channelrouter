package channelrouter

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

//ChannelRouter type:
type ChannelRouter struct {
	ingress  chan Packet
	run      bool
	running  bool
	Channels map[Key]*ChanLink
	wg       sync.WaitGroup

	Logger *log.Logger
}

// AddChannel adds a new underlying channel. A desired buffer size for the channel
// can be passed here. (default 32)
func (cr *ChannelRouter) AddChannel(buffer ...int) Key {
	var b int
	if len(buffer) < 1 || buffer[0] < 32 {
		b = 32
	} else {
		b = buffer[0]
	}
	cl := NewChanLink(b)
	cr.Channels[cl.Key] = cl
	cr.vog("Adding new channel to ChannelRouter with key:%v\n", cl.Key)
	return cl.Key
}

// SetType sets the type desired for the underlying channel referenced by k.
func (cr *ChannelRouter) SetType(k Key, t interface{}) {
	cr.vog("Setting type:%T for channel:%v\n", t, k)
	cr.Channels[k].refType = reflect.TypeOf(t)
}

// GetType returns the current type set for the underlying channel referenced by k.
func (cr *ChannelRouter) GetType(k Key) string {
	return fmt.Sprint(cr.Channels[k].refType)
}

// Stop stops the ChannelRouter process.
func (cr *ChannelRouter) Stop() {
	cr.vog("Issuing Stop() to ChannelRouter\n")
	cr.run = false
}

// Send sends i to the underlying channel referenced by k.
// Send will drop packets if the value of i does not match the value
// set for channel k unless the channel type is set to nil.
func (cr *ChannelRouter) Send(k Key, i interface{}) {
	cr.vog("Send request recieved")
	if cr.Channels[k].refType != nil {
		if cr.Channels[k].testTypeEq(i) != true {
			cr.vog("Dropping packet of type:%T >> channel:%v has specified type:%v\n", i, k, cr.Channels[k].refType)
			return
		}
	}
	//cr.wg.Add(1)
	p := Packet{
		header: k,
		value:  i,
	}
	cr.vog("Sending packet to ChannelRouter, destination:%v\n", k)
	cr.ingress <- p
	//cr.wg.Wait()
	cr.vog("%v", cr.statCheck(p.header))
}

// Receive returns a Packet from the underlying channel referenced by k.
// If no Packet is available, an empty Packet is returned with the empty packet
// error.
func (cr *ChannelRouter) Receive(k Key) Packet {
	cr.vog("Packet return requested for channel:%v\n", k)
	r, ok := cr.Channels[k].Receive()
	p := Packet{
		header: k,
	}
	if !ok {
		//add retry:
		//cr.vog("Nothing available to recieve, retyring once")
		//time.Sleep(time.Millisecond * 300)
		r, ok = cr.Channels[k].Receive()
		if !ok {
			cr.vog("Nothing available to recieve, returning nil")
			p.Err = fmt.Errorf("empty packet: nothing available in channel, ")
			p.value = nil
			return p
		}
	}
	p.value = r
	cr.Channels[k].received++
	cr.vog("Packet returned for channel:%v\n", k)
	return p
}

//Route starts the ChannelRouter process.
func (cr *ChannelRouter) Route() {
	cr.run = true
	go func() {
		for {
			if cr.run {
				select {
				case p := <-cr.ingress:
					cr.vog("ChannelRouter received packet, routing to channel:%v\n", p.header)
					if cr.Channels[p.header].sent > 1024 {
						reset := cr.Channels[p.header].resetCounters()
						cr.vog("Channel:%v counters reset:%v\n", p.header, reset)
					}
					cr.Channels[p.header].Send(p.value)
					cr.Channels[p.header].sent++
					cr.vog("Packet sent\n")
					//cr.wg.Done()
				default:
					cr.running = true
				}
				if !cr.run {
					return
				}
			} else {
				cr.vog("ChannelRouter no longer running!\n")
				break
			}
		}
	}()
	for {
		if cr.running {
			break
		}
	}
}

//Stats hold the packets Sent, Received and currently Available.
type Stats struct {
	Sent      uint32
	Received  uint32
	Available uint32
}

//GetStats returns the Stats type for the channel referenced by k.
func (cr *ChannelRouter) GetStats(k Key) Stats {
	return Stats{
		Sent:      cr.Channels[k].sent,
		Received:  cr.Channels[k].received,
		Available: cr.Channels[k].getAvailable(),
	}
}

func (cr *ChannelRouter) statCheck(k Key) string {
	sc := cr.GetStats(k)
	return fmt.Sprintf("Stats:%v >> Sent:%v Received:%v Available:%v", k, sc.Sent, sc.Received, sc.Available)
}

//Available returns the number of packets currently available in the channel specified by k.
func (cr *ChannelRouter) Available(k Key) uint32 {
	cr.vog("Getting number of outstanding packets available for channel:%v\n", k)
	return cr.Channels[k].getAvailable()
}

func (cr *ChannelRouter) listChannels() []Key {
	cr.vog("Getting all available channels.\n")
	var keys []Key
	for k := range cr.Channels {
		keys = append(keys, k)
	}
	return keys
}

func (cr *ChannelRouter) matchChannels(i interface{}) []Key {
	cr.vog("Matching any channels accepting type:%v\n", reflect.TypeOf(i))
	var keys []Key
	for _, k := range cr.listChannels() {
		if cr.Channels[k].testTypeEq(i) == true {
			keys = append(keys, k)
			cr.vog("Matched:%v\n", k)
		}
	}
	cr.vog("Found %v mathcing channels.\n", len(keys))
	return keys
}

// Broadcast sends to all underlying channels in ChannelRouter.
// If the underlying channel has a set type which doesn't match the broadcast type,
// the value will not be sent to that channel. If the underlying channel type has
// not been set, or set to nil, the value will be passed.
func (cr *ChannelRouter) Broadcast(i interface{}) {
	keys := cr.matchChannels(i)
	if len(keys) < 1 {
		cr.vog("No channels found accepting type:%v\n", reflect.TypeOf(i))
		return
	}
	cr.vog("Broadcast sending to all channels accepting type:%v\n", reflect.TypeOf(i))
	for _, k := range keys {
		cr.Send(k, i)
		cr.vog("Sent to channel:%v\n", k)
	}
}

// NewChannelRouter creates and returns a new ChannelRouter.
func NewChannelRouter(buffer ...int) *ChannelRouter {
	var b int
	if len(buffer) < 1 || buffer[0] < 1024 {
		b = 1024
	} else {
		b = buffer[0]
	}
	channels := make(map[Key]*ChanLink)
	return &ChannelRouter{
		ingress:  make(chan Packet, b),
		Channels: channels,
	}
}

func (cr *ChannelRouter) vog(f string, a ...interface{}) {
	if cr.Logger != nil {
		cr.Logger.Printf(f, a...)
	}
}
