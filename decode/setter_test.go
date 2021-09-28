package decode

import (
	"reflect"
	"testing"
)

type testType bool

func (t *testType) UnmarshalBin(b []byte) error {
	*t = b[0] & 0x1 != 0
	return nil
}

func getValue(fieldNum int, v interface{}) reflect.Value {
	return reflect.ValueOf(v).Elem().Field(fieldNum)
}

func TestSetValue(t *testing.T) {
	t.Run("negative scenario: user defined type without size. Expect error", func(t *testing.T) {
		var a struct{
			B testType
		}
		err := setValue(&reader{payload: []byte{1,2,3}},  &tag{
			arraySize:    nil,
			varSize:      nil,
			littleEndian: false,
		}, getValue(0, &a))

		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("negative scenario: user defined type - read error. Expect error", func(t *testing.T) {
		var a struct{
			B testType
		}
		size := int64(1)
		err := setValue(&reader{payload: []byte{}},  &tag{
			arraySize:    nil,
			varSize:      &size,
			littleEndian: false,
		}, getValue(0, &a))

		if err == nil {
			t.Error("expected error")
		}
	})
}