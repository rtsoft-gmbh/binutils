package decode

import (
	"reflect"
)

// Unmarshaler is the interface implemented by types
// that can unmarshal binary description of themselves.
type Unmarshaler interface {
	UnmarshalBin([]byte) error
}

// UnmarshalBin function used for binary unmarshal.
func UnmarshalBin(payload []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &ErrorInvalidUnmarshal{reflect.TypeOf(v)}
	}

	rd := &reader{
		pos:     0,
		payload: payload,
	}

	return unmarshalBin(rd, rv)
}

func unmarshalBin(rd *reader, rv reflect.Value) error {
	structValue := rv.Elem()
	structType := structValue.Type()
	for i := 0; i < structValue.NumField(); i++ {
		tag, err := parseTag(structType.Field(i))
		if err != nil {
			return err
		}
		err = setValue(rd, tag, structValue.Field(i))
		if err != nil {
			return err
		}
	}

	return nil
}

// ErrorInvalidUnmarshal can be returned if non pointer
// given in UnmarshalBin function
type ErrorInvalidUnmarshal struct {
	Type reflect.Type
}

func (e *ErrorInvalidUnmarshal) Error() string {
	if e.Type == nil {
		return "lora-payload: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Ptr {
		return "lora-payload: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "lora-payload: Unmarshal(nil " + e.Type.String() + ")"
}
