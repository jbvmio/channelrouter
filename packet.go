package channelrouter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

//Packet Here.
type Packet struct {
	err    error
	header Key
	tag    reflect.Type
	value  interface{}
}

/* Additional Packet Ideas:
// callback func
*/

//IsArray Here.
func (p Packet) IsArray() bool {
	return isArray(p.value)
}

func (p Packet) String() string {
	return cast.ToString(p.value)
}

//Int asserts int
func (p Packet) Int() int {
	return cast.ToInt(p.value)
}

//Array Here
func (p Packet) Array() ([]Packet, error) {
	if !isArray(p.value) {
		fmt.Println("")
		return []Packet{}, fmt.Errorf("this Packet does not contain an array")
	}
	key := p.header
	var packets []Packet
	s := reflect.ValueOf(p.value)

	for i := 0; i < s.Len(); i++ {
		pkt := Packet{
			header: key,
			value:  s.Index(i),
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func isArray(i interface{}) bool {
	r := reflect.ValueOf(i)
	typeOfT := r.Type()
	t := fmt.Sprint(typeOfT)
	return strings.Contains(t, "[]")
}
