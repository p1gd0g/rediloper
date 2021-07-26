package example

import (
	"log"
	"reflect"
	"testing"

	"github.com/p1gd0g/rediloper/example/rawproto"
)

func TestNewPerson(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name string
		args args
		want *Person
	}{
		{
			name: "",
			args: args{
				uid: "1000",
			},
			want: &Person{
				Meta: &Meta{
					params: []interface{}{"1000"},
					format: "Person:%s",
				},
				Message: &rawproto.Person{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPerson(tt.args.uid)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPerson() = %v, want %v", got, tt.want)
			}
			log.Println("The key is", got.GetKey())
			got.GetObj().Name = tt.args.uid
			b, err := got.Marshal()
			if err != nil {
				t.Error(err)
			}
			cp := NewPerson(tt.args.uid)
			err = cp.Unmarshal(b)
			if err != nil {
				t.Error(err)
			}
			log.Println(got.Message, cp.Message)
		})
	}

}
