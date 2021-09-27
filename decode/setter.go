package decode

import (
	"encoding"
	"encoding/binary"
	"errors"
	"math"
	"reflect"
)

func setValue(r *reader, t *tag, fieldValue reflect.Value) error {
	u, ut, pv := indirect(fieldValue, false)
	if u != nil {
		if t.varSize == nil {
			return errors.New("var_size should be specified for defined handlers")
		}
		n := *t.varSize
		pl, err := r.getBytes(int(n))
		if err != nil {
			return err
		}
		return u.UnmarshalBin(pl)
	}
	if ut != nil {
		return errors.New("type error: " + fieldValue.Kind().String())
	}

	fieldValue = pv

	var byteOrder binary.ByteOrder
	if t.littleEndian {
		byteOrder = binary.LittleEndian
	} else {
		byteOrder = binary.BigEndian
	}
	// template taken from: https://github.com/ghostiam/binstruct
	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var value int64
		var e error

		switch {
		case t.varSize != nil && *t.varSize == 1 || fieldValue.Kind() == reflect.Int8:
			pl, err := r.getBytes(1)
			if err == nil {
				value = int64(pl[0])
			}
			e = err

		case t.varSize != nil && *t.varSize == 2 || fieldValue.Kind() == reflect.Int16:
			pl, err := r.getBytes(2)
			if err == nil {
				value = int64(byteOrder.Uint16(pl))
			}
			e = err
		case t.varSize != nil && *t.varSize == 4 || fieldValue.Kind() == reflect.Int32:
			pl, err := r.getBytes(4)
			if err == nil {
				value = int64(byteOrder.Uint32(pl))
			}
			e = err
		case t.varSize != nil && *t.varSize == 8 || fieldValue.Kind() == reflect.Int64:
			pl, err := r.getBytes(8)
			if err == nil {
				value = int64(byteOrder.Uint64(pl))
			}
			e = err
		default: // reflect.Int:
			e = errors.New("need set tag with var_size or use int8/int16/int32/int64")
		}
		if e != nil {
			return e
		}
		if fieldValue.CanSet() {
			fieldValue.SetInt(value)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var value uint64
		var e error

		switch {
		case t.varSize != nil && *t.varSize == 1 || fieldValue.Kind() == reflect.Uint8:
			pl, err := r.getBytes(1)
			if err == nil {
				value = uint64(pl[0])
			}
			e = err

		case t.varSize != nil && *t.varSize == 2 || fieldValue.Kind() == reflect.Uint16:
			pl, err := r.getBytes(2)
			if err == nil {
				value = uint64(byteOrder.Uint16(pl))
			}
			e = err
		case t.varSize != nil && *t.varSize == 4 || fieldValue.Kind() == reflect.Uint32:
			pl, err := r.getBytes(4)
			if err == nil {
				value = uint64(byteOrder.Uint32(pl))
			}
			e = err
		case t.varSize != nil && *t.varSize == 8 || fieldValue.Kind() == reflect.Uint64:
			pl, err := r.getBytes(8)
			if err == nil {
				value = byteOrder.Uint64(pl)
			}
			e = err
		default: // reflect.Uint:
			e = errors.New("need set tag with varSize or use int8/int16/int32/int64")
		}
		if e != nil {
			return e
		}
		if fieldValue.CanSet() {
			fieldValue.SetUint(value)
		}

	case reflect.Float32:
		pl, err := r.getBytes(4)
		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetFloat(float64(math.Float32frombits(byteOrder.Uint32(pl))))
		}

	case reflect.Float64:
		pl, err := r.getBytes(8)
		if err != nil {
			return err
		}
		if fieldValue.CanSet() {
			fieldValue.SetFloat(math.Float64frombits(byteOrder.Uint64(pl)))
		}

	case reflect.Bool:
		b, err := r.getBytes(1)
		if err != nil {
			return err
		}
		if fieldValue.CanSet() {
			fieldValue.SetBool(b[0] == 1)
		}

	case reflect.String:
		if t.arraySize == nil {
			return errors.New("need set tag with array_size for string")
		}

		b, err := r.getBytes(int(*t.arraySize))
		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetString(string(b))
		}

	case reflect.Slice:
		if t.arraySize == nil {
			return errors.New("need set tag with array_size for slice")
		}

		for i := int64(0); i < *t.arraySize; i++ {
			tmpV := reflect.New(fieldValue.Type().Elem()).Elem()
			err := setValue(r, &tag{varSize: t.varSize, littleEndian: t.littleEndian}, tmpV)
			if err != nil {
				return err
			}
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.Append(fieldValue, tmpV))
			}
		}

	case reflect.Array:
		var arrLen int64

		if t.arraySize != nil {
			arrLen = *t.arraySize
		}

		if arrLen == 0 {
			arrLen = int64(fieldValue.Len())
		}

		for i := int64(0); i < arrLen; i++ {
			tmpV := reflect.New(fieldValue.Type().Elem()).Elem()
			err := setValue(r, &tag{varSize: t.varSize, littleEndian: t.littleEndian}, tmpV)
			if err != nil {
				return err
			}
			if fieldValue.CanSet() {
				fieldValue.Index(int(i)).Set(tmpV)
			}
		}

	case reflect.Struct:
		err := unmarshalBin(r, reflect.ValueOf(fieldValue.Addr().Interface()))
		if err != nil {
			return err
		}

	default:
		return errors.New(`type "` + fieldValue.Kind().String() + `" not supported`)
	}

	return nil
}

func indirect(v reflect.Value, decodingNull bool) (Unmarshaler, encoding.TextUnmarshaler, reflect.Value) {
	v0 := v
	haveAddr := false

	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if decodingNull && v.CanSet() {
			break
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		//     var v interface{}
		//     v = &v
		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 && v.CanInterface() {
			if u, ok := v.Interface().(Unmarshaler); ok {
				return u, nil, reflect.Value{}
			}
			if !decodingNull {
				if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
					return nil, u, reflect.Value{}
				}
			}
		}

		if haveAddr {
			v = v0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}
	return nil, nil, v
}
