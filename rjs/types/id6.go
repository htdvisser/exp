package types

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

type ID6 [8]byte

func (id ID6) MarshalText() ([]byte, error) {
	var res bytes.Buffer
	if id[0] != 0 || id[1] != 0 {
		res.WriteString(fmt.Sprintf("%x:", uint16(id[0])<<8|uint16(id[1])))
	}
	for _, g := range []uint16{
		uint16(id[2])<<8 | uint16(id[3]),
		uint16(id[4])<<8 | uint16(id[5]),
	} {
		if g != 0 {
			res.WriteString(fmt.Sprintf("%x", g))
		}
		res.WriteString(":")
	}
	res.WriteString(fmt.Sprintf("%x", uint16(id[6])<<8|uint16(id[7])))
	return res.Bytes(), nil
}

var emptyID6 ID6

func (id ID6) String() string {
	if id == emptyID6 {
		return ""
	}
	b, _ := id.MarshalText()
	return string(b)
}

func (id ID6) EUIString() string {
	return fmt.Sprintf("%X", id[:])
}

var (
	euiPattern = regexp.MustCompile(`^([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})$`)
	macPattern = regexp.MustCompile(`^([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})[-:]?([a-fA-F0-9]{2})$`)
	id6Pattern = regexp.MustCompile(`^(?:([a-f0-9]{0,4}):)?([a-f0-9]{0,4}):([a-f0-9]{0,4}):([a-f0-9]{0,4})$`)
)

func (id *ID6) UnmarshalText(data []byte) error {
	if bytes := euiPattern.FindStringSubmatch(string(data)); bytes != nil {
		for i, b := range bytes[1:] {
			v, err := strconv.ParseUint(b, 16, 8)
			if err != nil {
				return fmt.Errorf("invalid data: %w", err)
			}
			id[i] = uint8(v)
		}
		return nil
	}
	if bytes := macPattern.FindStringSubmatch(string(data)); bytes != nil {
		for i, b := range bytes[1:] {
			v, err := strconv.ParseUint(b, 16, 8)
			if err != nil {
				return fmt.Errorf("invalid data: %w", err)
			}
			if i < 3 {
				id[i] = uint8(v)
			} else {
				id[i+2] = uint8(v)
			}
		}
		id[3], id[4] = 0xff, 0xfe
		return nil
	}
	if bytes := id6Pattern.FindStringSubmatch(string(data)); bytes != nil {
		for i, b := range bytes[1:] {
			if b == "" {
				b = "0"
			}
			v, err := strconv.ParseUint(b, 16, 16)
			if err != nil {
				return fmt.Errorf("invalid data: %w", err)
			}
			id[2*i] = uint8(v >> 8)
			id[2*i+1] = uint8(v)
		}
		return nil
	}
	return fmt.Errorf("invalid data")
}
