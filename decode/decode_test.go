package decode

import (
	"testing"
)

func TestUnmarshalBin(t *testing.T) {
	t.Run("unmarshal nil. Expect error", func(t *testing.T) {
		err := UnmarshalBin([]byte{1,2,3}, nil)
		if err == nil {
			t.Error("expected error")
		}
		_ = err.Error()
	})
	
	t.Run("unmarshal non-pointer. Expect error", func(t *testing.T) {
		var a struct {}
		err := UnmarshalBin([]byte{1,2,3}, a)
		if err == nil {
			t.Error("expected error")
		}
		_ = err.Error()
	})

	t.Run("unmarshal nil pointer. Expect error", func(t *testing.T) {
		var a struct{}
		k := &a
		k = nil
		err := UnmarshalBin([]byte{1,2,3}, k)
		if err == nil {
			t.Error("expected error")
		}
		_ = err.Error()
	})

	t.Run("positive scenario", func(t *testing.T) {
		var a struct{
			Field int8
		}
		err := UnmarshalBin([]byte{0xF}, &a)
		if err != nil {
			t.Error("unexpected error: " + err.Error())
		}
		if a.Field != 0xF {
			t.Error("bad unmarshalling")
		}
	})
	
	t.Run("negative scenario: bad tag. Expect error", func(t *testing.T) {
		var s struct {
			t int `var_size:"asdas"`
		}
		err := UnmarshalBin([]byte{0xF}, &s)
		if err == nil {
			t.Error("Expected error at bad tag")
		}
	})

	t.Run("negative scenario: int without size. Expect error", func(t *testing.T) {
		var s struct{
			t int
		}
		err := UnmarshalBin([]byte{0xF}, &s)
		if err == nil {
			t.Error("Expected error at no size specified")
		}
	})
}
