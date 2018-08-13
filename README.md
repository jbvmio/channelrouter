# ChannelRouter: Package for Go

The ChannelRouter package enables quick creation and use of channels. Channels created with ChannelRouter are referenced by a unique Key which is used to "route" values to an appropriate underlying channel either directly or by broadcast.

The goal of ChannelRouter is quick and easy channel setup, flexibility and additional features for increased functionality. 

# Quick Example:

```
package main

import (
	"fmt"

	"github.com/jbvmio/channelrouter"
)

func main() {
	cr := channelrouter.NewChannelRouter()
	key1 := cr.AddChannel(32)
	cr.Route()

	cr.Send(key1, 777)
	a := cr.Receive(key1)
	fmt.Println(a)
}
```
*Quick Example Output:*
```
777
```

# Set Channel Type:
For some added control, the channel type can be set enabling ChannelRouter to allow only the specified type to be sent. This is done by simply passing an empty value of the type you want to set on the channel.
```
package main

import (
	"fmt"
	"time"

	"github.com/jbvmio/channelrouter"
)

func main() {
	cr := channelrouter.NewChannelRouter()
	//cr.Logger = log.New(os.Stdout, "[channelRouter] ", log.LstdFlags)
	key1 := cr.AddChannel(32)
	cr.Route()

	cr.SetType(key1, "")    //set channel type to string

	i := 777
	cr.Send(key1, i)

	s := "777"
	cr.Send(key1, s)

	for cr.Available(key1) < 1 {
		time.Sleep(time.Second * 1)
	}
	for cr.Available(key1) > 0 {
		a := cr.Receive(key1)
		fmt.Printf("Value: %v Type: %T\n", a, a.Value())
	}
}
```
*Set Channel Output:*
```
Value: 777 Type: string
```
If you uncomment the logger, you can see the first value of 777 (int) was dropped:
```
...
[channelRouter] 2018/08/13 14:59:25 Setting type:string for channel:5e1a0c91da9ea5c9
[channelRouter] 2018/08/13 14:59:25 Send request recieved
[channelRouter] 2018/08/13 14:59:25 Dropping packet of type:int >> channel:5e1a0c91da9ea5c9 has specified type:string
...
Value: 777 Type: string
[channelRouter] 2018/08/13 14:59:26 Getting number of outstanding packets available for channel:5e1a0c91da9ea5c9
...
```

# Broadcasting, adding multiple channels:
If you use ChannelRouter to send a broadcast, the type sent will be matched to all other channels that are set to the same type, as well as channels that have not been set at all, ie. <nil>
```
package main

import (
	"fmt"

	"github.com/jbvmio/channelrouter"
)

func main() {
	cr := channelrouter.NewChannelRouter()
	intChan := cr.AddChannel(32)                  //Add the first channel
	stringChan := cr.AddChannel(32)               //Add the second channel
	cr.Route()

	cr.SetType(intChan, int(0))           //set channel type to int
	cr.SetType(stringChan, string(""))    //set channel type to string

	cr.Broadcast(777)
	fmt.Printf("Available values in intChan: %v\n", cr.Available(intChan))        //Here
	fmt.Printf("Available values in stringChan: %v\n", cr.Available(stringChan))  //Not Here

	cr.Broadcast("777")
	fmt.Printf("Available values in intChan: %v\n", cr.Available(intChan))        //Not Here
	fmt.Printf("Available values in stringChan: %v\n", cr.Available(stringChan))  //Here

	i := cr.Receive(intChan)
	s := cr.Receive(stringChan)
	fmt.Printf("Value:%v Type:%T\n", i, i.Value())
	fmt.Printf("Value:%v Type:%T\n", s, s.Value())
}

```
*Output:*
```
Available values in intChan: 1
Available values in stringChan: 0
Available values in intChan: 1
Available values in stringChan: 1
Value:777 Type:int
Value:777 Type:string
```

# Use as an io.Writer:
ChannelRouter can also be used as an io.Writer:
```
package main

import (
	"fmt"

	"github.com/jbvmio/channelrouter"
)

func main() {
	cr := channelrouter.NewChannelRouter(1024)
	byteChan := cr.AddChannel(32)
	cr.Route()

	io := cr.MakeIoChannel(byteChan)

	fmt.Fprintln(io, "Hello There!")

	var line []byte
	for cr.Available(byteChan) > 0 {
		a := cr.Receive(byteChan).ToByte()
		line = append(line, a)
	}
	fmt.Printf(string(line))
}

```
*io.Writer Output:*
```
Hello There!
```

**Note:**
* This package is currently under development
