package channelrouter

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

//ChannelRouter Here.
type ChannelRouter struct {
	ingress  chan Packet
	run      bool
	running  bool
	Channels map[Key]*ChanLink
	wg       sync.WaitGroup

	Logger *log.Logger
}

//AddChannel Here.
func (cr *ChannelRouter) AddChannel(buffer int) Key {
	cl := NewChanLink(buffer)
	cr.Channels[cl.Key] = cl
	cr.vog("Adding new channel to ChannelRouter with key:%v\n", cl.Key)
	return cl.Key
}

//SetType Here.
func (cr *ChannelRouter) SetType(k Key, t interface{}) {
	cr.vog("Setting type:%T for channel:%v\n", t, k)
	cr.Channels[k].refType = reflect.TypeOf(t)
}

//Stop Here.
func (cr *ChannelRouter) Stop() {
	cr.vog("Issuing Stop() to ChannelRouter\n")
	cr.run = false
}

//Send Here.
func (cr *ChannelRouter) Send(k Key, i interface{}) {
	cr.vog("Send request recieved")
	if cr.Channels[k].refType != nil {
		if cr.Channels[k].testTypeEq(i) != true {
			cr.vog("Dropping packet of type:%T >> channel:%v has specified type:%v\n", i, k, cr.Channels[k].refType)
			return
		}
	}
	cr.wg.Add(1)
	p := Packet{
		header: k,
		value:  i,
	}
	cr.vog("Sending packet to ChannelRouter, destination:%v\n", k)
	cr.ingress <- p
	cr.wg.Wait()
	cr.vog("%v", cr.statCheck(p.header))
}

//Receive Here.
func (cr *ChannelRouter) Receive(k Key) Packet {
	cr.vog("Packet return requested for channel:%v\n", k)
	r, ok := cr.Channels[k].Receive()
	p := Packet{
		header: k,
	}
	if !ok {
		cr.vog("Nothing available to recieve, returning nil")
		p.value = nil
		return p
	}
	p.value = r
	cr.Channels[k].received++
	cr.vog("Packet returned for channel:%v\n", k)
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
					cr.vog("ChannelRouter received packet, routing to channel:%v\n", p.header)
					if cr.Channels[p.header].sent > 1024 {
						reset := cr.Channels[p.header].resetCounters()
						cr.vog("Channel:%v counters reset:%v\n", p.header, reset)
					}
					cr.Channels[p.header].Send(p.value)
					cr.Channels[p.header].sent++
					cr.vog("Packet sent\n")
					cr.wg.Done()
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

//Stats Here.
type Stats struct {
	Sent      uint32
	Received  uint32
	Available uint32
}

//GetStats Here.
func (cr *ChannelRouter) GetStats(k Key) Stats {
	return Stats{
		Sent:      cr.Channels[k].sent,
		Received:  cr.Channels[k].received,
		Available: cr.Channels[k].getAvailable(),
	}
}

func (cr *ChannelRouter) statCheck(k Key) string {
	sc := cr.GetStats(k)
	return fmt.Sprintf("Status Check >> Sent:%v Received:%v Available:%v", sc.Sent, sc.Received, sc.Available)
}

//Available Here.
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

//Broadcast Here.
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

//NewChannelRouter Here.
func NewChannelRouter() *ChannelRouter {
	channels := make(map[Key]*ChanLink)
	return &ChannelRouter{
		ingress:  make(chan Packet, 1024),
		Channels: channels,
	}
}

//Stop Here.
func (cr *ChannelRouter) vog(f string, a ...interface{}) {
	if cr.Logger != nil {
		cr.Logger.Printf(f, a...)
	}
}
