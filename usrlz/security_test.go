package usrlz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type securityPaddedValue struct {
	Prefix uint8
	Number uint32
}

func TestToBytesUsesDeterministicLittleEndianEncoding(t *testing.T) {
	value := uint32(0x01020304)
	assert.Equal(t, []byte{0x04, 0x03, 0x02, 0x01}, ToBytes(&value))
}

func TestToBytesEncodesPlatformIntegersAs64BitValues(t *testing.T) {
	value := int(0x01020304)
	assert.Equal(t, []byte{0x04, 0x03, 0x02, 0x01, 0, 0, 0, 0}, ToBytes(&value))
}

func TestToBytesOmitsStructPadding(t *testing.T) {
	value := securityPaddedValue{Prefix: 0x11, Number: 0x22334455}
	assert.Equal(t, []byte{0x11, 0x55, 0x44, 0x33, 0x22}, ToBytes(&value))
}

func TestToBytesRejectsNilPointers(t *testing.T) {
	var value *uint32
	assert.PanicsWithValue(t, "object is nil", func() {
		ToBytes(value)
	})
}

func TestToBytesRejectsReferenceBearingValues(t *testing.T) {
	value := struct {
		Text string
	}{Text: "secret"}

	assert.Panics(t, func() {
		ToBytes(&value)
	})
}
