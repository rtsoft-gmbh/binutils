package decode

import (
	"reflect"
	"testing"
)

type arraySizeTestOK struct {
	t int `array_size:"10"`
}

type arraySizeTestBad struct {
	t int `array_size:"asdas"`
}

func TestParseTag(t *testing.T) {
	t.Run("positive scenario. array_size", func(t *testing.T) {
		var s arraySizeTestOK
		tg, err := parseTag(reflect.ValueOf(&s).Elem().Type().Field(0))
		if err != nil {
			t.Error("Unexpected error: " + err.Error())
		}
		if tg.arraySize == nil || *tg.arraySize != 10 {
			t.Error("Unexpected tag")
		}
	})
	t.Run("negative scenario. Cannot parse array_size. Expect error", func(t *testing.T) {
		var s arraySizeTestBad
		_, err := parseTag(reflect.ValueOf(&s).Elem().Type().Field(0))
		if err == nil {
			t.Error("Expected error")
		}
	})
}
