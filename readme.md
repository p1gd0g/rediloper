# Rediloper

`rediloper` is a protobuf wrapper generator, designed for `mget` or `mset` protobuf variables easily.

## Installation

```
go get -u github.com/p1gd0g/rediloper
go build
```

## Quick start

Given a protobuf file:

```protobuf
message Person {
    string name  = 1;
    int32  id    = 2;
    string email = 3;
}
```

Add redis key info like `@PE|name|format|param1|type1(|param2|type2|...)|comment`:

```protobuf
//@PE|Person|Person:%s|uid|string|person data
message Person {
    string name  = 1;
    int32  id    = 2;
    string email = 3;
}
```

Generate protobuf go file:

```
$ protoc -I=./example --go_out=. ./example/addressbook.proto
```

Generate rediloper wrapper:

```
./rediloper example/rawproto/addressbook.pb.go ./example/ ./
```

The output would be:

(or edit templates/file.tpl to match your case)

```go
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
```

## So, how to mget or mset?

I can't give a general solution for mget or mset, it depends on which redis client you are using and your redis conn pool.

But it would not be far away from the following:

1. Implement the `GetKey` method for `Meta` type:

```go
func (m *Meta) GetKey() (s string) {
	stringParams := make([]string, len(m.params))
	for i, param := range m.params {
		switch param := param.(type) {
		case string:
			stringParams[i] = param
		default:
			stringParams[i] = fmt.Sprint(param)
		}
	}

	s = strings.Join(stringParams, ":")
	s = fmt.Sprintf(m.format, s)
	return
}
```

2. Code the `MGET` function:

```go
func MGET(wrappers ...ProtoWrapper) error {

	cli := &rediscli{}

	keys := make([]string, 0, len(wrappers))
	for _, v := range wrappers {
		keys = append(keys, v.GetKey())
	}

	result, err := cli.MGET(keys)
	if err != nil {
		return err
	}

	for i, v := range result {
		err = wrappers[i].Unmarshal(v)
		if err != nil {
			return err
		}
	}

	return nil
}
```

See example for details.