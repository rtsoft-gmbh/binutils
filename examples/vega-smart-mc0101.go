package main

import (
	"errors"
	"fmt"

	"github.com/rtsoft-gmbh/binutils/decode"
)

type pkgType uint8

// packages types
const (
	DefaultPackage pkgType = iota
	EventPackage
)

var packageTypes = map[byte]pkgType{
	0x01: DefaultPackage,
	0xFF: EventPackage,
}

func getPackageType(payload []byte) (*pkgType, error) {
	if len(payload) < 1 {
		return nil, errors.New("bad package")
	}
	if v, ok := packageTypes[payload[0]]; ok {
		return &v, nil
	}
	return nil, errors.New("package this type not found")
}

func (p *pkgType) UnmarshalBin(b []byte) error {
	pt, err := getPackageType(b)
	if err != nil {
		return err
	}
	*p = *pt

	return nil
}

type pinStatuses struct {
	OpenedPin1 bool
	OpenedPin2 bool
}

func (p *pinStatuses) UnmarshalBin(b []byte) error {
	p.OpenedPin1 = b[0]&0x01 == 0
	p.OpenedPin2 = b[0]&0x02 == 0

	return nil
}

type sendCause uint8

const (
	timeCause    sendCause = 0
	sensor1Event           = 1
	sensor2Event           = 2
)

func (s *sendCause) UnmarshalBin(b []byte) error {
	if b[0] > 2 {
		return errors.New("bad value")
	}
	*s = sendCause(b[0])
	return nil
}

type vegaSmartMC0101DefaultPackage struct {
	PackageType pkgType `var_size:"1"`
	Battery     byte
	NullByte    byte
	Temperature int16       `byte_order:"le"`
	SendCause   sendCause   `var_size:"1"`
	PinStatuses pinStatuses `var_size:"1"`
	UnixTime    uint32      `byte_order:"le"`
}

func processMessageVEGASmartMC0101(payload []byte) error {
	pkg, err := getPackageType(payload)
	if err != nil {
		return err
	}
	if *pkg != DefaultPackage {
		return errors.New("only default package supported yet")
	}
	var vegaPackage vegaSmartMC0101DefaultPackage
	err = decode.UnmarshalBin(payload, &vegaPackage)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", vegaPackage)
	return nil
}

func main() {
	err := processMessageVEGASmartMC0101([]byte{1, 90, 0, 160, 0, 0, 2, 54, 179, 81, 97})
	fmt.Println(err)
}
