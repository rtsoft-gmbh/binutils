package decode

import (
	"reflect"
	"strconv"
)

type tag struct {
	arraySize    *int64
	varSize      *int64
	littleEndian bool
}

func parseTag(sf reflect.StructField) (*tag, error) {
	tg := &tag{}
	st := sf.Tag
	if val, ok := st.Lookup("array_size"); ok {
		len, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		tg.arraySize = &len
	}
	if val, ok := st.Lookup("var_size"); ok {
		len, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		tg.varSize = &len
	}
	if val, ok := st.Lookup("byte_order"); ok && val == "le" {
		tg.littleEndian = true
	}
	return tg, nil
}
