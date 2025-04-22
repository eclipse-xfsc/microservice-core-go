package ptr_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/ptr"
)

func TestBool(t *testing.T) {
	assert.True(t, *ptr.Bool(true))
	assert.False(t, *ptr.Bool(false))
}

func TestByte(t *testing.T) {
	assert.Equal(t, byte(64), *ptr.Byte(64))
}

func TestInt(t *testing.T) {
	assert.Equal(t, int(64), *ptr.Int(64))
}

func TestInt8(t *testing.T) {
	assert.Equal(t, int8(64), *ptr.Int8(64))
}

func TestInt16(t *testing.T) {
	assert.Equal(t, int16(64), *ptr.Int16(64))
}

func TestInt32(t *testing.T) {
	assert.Equal(t, int32(64), *ptr.Int32(64))
}

func TestInt64(t *testing.T) {
	assert.Equal(t, int64(64), *ptr.Int64(64))
}

func TestUint(t *testing.T) {
	assert.Equal(t, uint(64), *ptr.Uint(64))
}

func TestUint8(t *testing.T) {
	assert.Equal(t, uint8(64), *ptr.Uint8(64))
}

func TestUint16(t *testing.T) {
	assert.Equal(t, uint16(64), *ptr.Uint16(64))
}

func TestUint32(t *testing.T) {
	assert.Equal(t, uint32(64), *ptr.Uint32(64))
}

func TestUint64(t *testing.T) {
	assert.Equal(t, uint64(64), *ptr.Uint64(64))
}

func TestFloat32(t *testing.T) {
	assert.Equal(t, float32(42.5), *ptr.Float32(42.5))
}

func TestFloat64(t *testing.T) {
	assert.Equal(t, float64(42.5), *ptr.Float64(42.5))
}

func TestRune(t *testing.T) {
	assert.Equal(t, 'r', *ptr.Rune('r'))
}

func TestString(t *testing.T) {
	assert.Equal(t, "mystring", *ptr.String("mystring"))
}

func TestTime(t *testing.T) {
	now := time.Now()
	assert.Equal(t, now, *ptr.Time(now))
}
