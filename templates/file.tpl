{{ define "file" }}

package example

import (
   	"github.com/p1gd0g/rediloper/example/rawproto"
	"google.golang.org/protobuf/proto"
)

{{ range .SheetConfigs.PE }}
type {{ .ClassName }} struct {
        *Meta
        proto.Message
}

func (p *{{ .ClassName }}) GetProto() proto.Message {
	return p.Message
}

func (p *{{ .ClassName }}) Marshal() ([]byte, error) {
	return proto.Marshal(p.Message)
}

func (p *{{ .ClassName }}) Unmarshal(b []byte) error {
	return proto.Unmarshal(b, p.Message)
}

func (p *{{ .ClassName }}) GetObj() *rawproto.{{ .TypeName }} {
	return p.Message.(*rawproto.{{ .TypeName }})
}

func New{{ .ClassName }}({{ range .KeyInputs }}{{ . }},{{ end }}) *{{ .ClassName }} {
    return &{{ .ClassName }}{Message: &rawproto.{{ .TypeName }}{}, Meta: &Meta{format: "{{.KeyFormat}}", params: []interface{}{ {{ range .KeyParams }} {{ . }}, {{ end }} }}}
}

{{ end }}

{{ end }}