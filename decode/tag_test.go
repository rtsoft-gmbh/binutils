package decode

import (
	"reflect"
	"testing"
)

func TestParseTag(t *testing.T) {
	t.Run("positive scenario: array_size", func(t *testing.T) {
		var s struct {
			t int `array_size:"10"`
			k int `var_size:"10"`
			e int `byte_order:"le"`
		}
		tg, err := parseTag(reflect.ValueOf(&s).Elem().Type().Field(0))
		if err != nil {
			t.Error("Unexpected error: " + err.Error())
		}
		if tg.arraySize == nil || *tg.arraySize != 10 {
			t.Errorf("Unexpected tag. tg %+v", tg)
		}
	})

	t.Run("positive scenario: var_size", func(t *testing.T) {
		var s struct {
			k int `var_size:"10"`
		}
		tg, err := parseTag(reflect.ValueOf(&s).Elem().Type().Field(0))
		if err != nil {
			t.Error("Unexpected error: " + err.Error())
		}
		if tg.varSize == nil || *tg.varSize != 10 {
			t.Errorf("Unexpected tag. tg %+v", tg)
		}
	})

	t.Run("positive scenario: byte_order", func(t *testing.T) {
		var s struct {
			k int `byte_order:"le"`
		}
		tg, err := parseTag(reflect.ValueOf(&s).Elem().Type().Field(0))
		if err != nil {
			t.Error("Unexpected error: " + err.Error())
		}
		if !tg.littleEndian {
			t.Errorf("Unexpected tag. tg %+v", tg)
		}
	})

	t.Run("negative scenario. Cannot parse array_size. Expect error", func(t *testing.T) {
		var s struct {
			t int `array_size:"asdas"`
		}
		_, err := parseTag(reflect.ValueOf(&s).Elem().Type().Field(0))
		if err == nil {
			t.Error("Expected error")
		}
	})
}
