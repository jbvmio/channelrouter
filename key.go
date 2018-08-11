package channelrouter

import "github.com/jbvmio/randstr"

//Key Here.
type Key string

func newKey() Key {
	t := randstr.Hex(8)
	return Key(t)
}
