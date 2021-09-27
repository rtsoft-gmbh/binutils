## Simple lib for binary decoding (encoding in feature)

### Description
This library can be used to decode binary payload in to the structures.  
The main idea was to make the library simple and fast.


### Use
Basic types can be used. (with tag) - means that array_size or var_size tags should be set:
* uint8, uint16, uint32, uint64. uint(with tag)
* int8. int16, int32, int64. int(with tag var_size)
* slices (with tag array_size \[& var_size\])
* arrays (with tag array_size \[& var_size\])
* user-defined types (with tag var_size)
* strings (with tag array_size)
* structs

```go
type anotherStruct struct {
	SomeValue int8 
	AbotherValue uint64 `var_size:"2"`
}

type someStruct struct {
	PackageType pkgType `var_size:"1"`
	Battery     byte
	NullByte    byte
	SendCause   sendCause   `var_size:"1"`
	PinStatuses pinStatuses `var_size:"1"`
	Message     string  `array_size:"10"`
	SomeElse    anotherStruct
}
```

#### Tags
You can use 3 type of tags:
* `array_size`: used to define count of the elements for slices, arrays(without defined length) and strings
* `var_size`: can be used with any type except string. Required for the int, uint and user defined types. It can be used to decode fewer bytes that variable requires. For example var_size: "1" for uint64 will read 1 byte and convert it to uint64. If `var_size` contains more bytes than variable requires, the tag will be ignored.
* `byte_order:"le"` used to set a byte order to the little endian. Default is the big endian.

```go
type temperaturePackage struct {
	Temperature int16  `byte_order:"le"`
	UnixTime    uint32  `byte_order:"le"`
}
```

#### Define unmarshaler
You can define your own unmarshaler for your types. It works the same as the UnmarshalJson in "encoding/json" package.  
Note: `var_size` tag should be set. `byte_order`,`array_size` not works for this case. The number of bytes set in var_size will be given in UnmarshalBin function.  
Example. 1 byte will be given in UnmarshalBin function:  
```go
type pinStatuses struct {
	OpenedPin1 bool
	OpenedPin2 bool
}

func (p *pinStatuses) UnmarshalBin(b []byte) error {
	p.OpenedPin1 = b[0]&0x01 == 0
	p.OpenedPin2 = b[0]&0x02 == 0

	return nil
}

type someStruct struct {
	PinStatuses pinStatuses `var_size:"1"`
}
```
### Features
* codecov
* docs
* examples
* encoder
### Examples
See tests and examples for more information.
### License
MIT License