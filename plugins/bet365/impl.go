package bet365

import (
	"fmt"

	"github.com/bet365/jingo"
)

type Payload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Bet365Plugin struct {
	enc *jingo.StructEncoder
}

func (b *Bet365Plugin) Initialize() {
	b.enc = jingo.NewStructEncoder(Payload{})
}

func (b *Bet365Plugin) CallMethod() string {

	payload := Payload{
		Name: "Who",
		Age:  33,
	}

	buffer := jingo.NewBufferFromPool()
	b.enc.Marshal(&payload, buffer)
	fmt.Println(buffer.String())
	buffer.ReturnToPool()

	return ""
}
