package example

import (
	"github.com/p1gd0g/rediloper/example/rawproto"
	"google.golang.org/protobuf/proto"
)

type Person struct {
	*Meta
	proto.Message
}

func (p *Person) GetProto() proto.Message {
	return p.Message
}

func (p *Person) Marshal() ([]byte, error) {
	return proto.Marshal(p.Message)
}

func (p *Person) Unmarshal(b []byte) error {
	return proto.Unmarshal(b, p.Message)
}

func (p *Person) GetObj() *rawproto.Person {
	return p.Message.(*rawproto.Person)
}

func NewPerson(uid string) *Person {
	return &Person{Message: &rawproto.Person{}, Meta: &Meta{format: "Person:%s", params: []interface{}{uid}}}
}
