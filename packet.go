package channelrouter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

//Packet Here.
type Packet struct {
	Err    error
	header Key
	tag    reflect.Type
	value  interface{}
}

/* Additional Packet Ideas:
// callback func
*/

//Error Here.
func (p Packet) Error() error {
	return p.Err
}

//Value Here.
func (p Packet) Value() interface{} {
	return p.value
}

//ToString assets string
func (p Packet) String() string {
	return cast.ToString(p.value)
}

//IsArray Here.
func (p Packet) IsArray() bool {
	return isArray(p.value)
}

//ToString assets string
func (p Packet) ToString() string {
	return cast.ToString(p.value)
}

//ToInt asserts int
func (p Packet) ToInt() int {
	return cast.ToInt(p.value)
}

//ToByte asserts int
func (p Packet) ToByte() byte {
	return cast.ToUint8(p.value)
}

//ToArray Here
func (p Packet) ToArray() ([]Packet, error) {
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
