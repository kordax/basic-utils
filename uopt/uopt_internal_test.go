package uopt

import (
	"database/sql/driver"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingValuer struct{}

func (failingValuer) Value() (driver.Value, error) {
	return nil, assert.AnError
}

func TestScanDuration(t *testing.T) {
	value, err := scanDuration(int64(time.Second))
	require.NoError(t, err)
	assert.Equal(t, time.Second, value)

	value, err = scanDuration(float64(time.Millisecond))
	require.NoError(t, err)
	assert.Equal(t, time.Millisecond, value)

	value, err = scanDuration("2s")
	require.NoError(t, err)
	assert.Equal(t, 2*time.Second, value)

	value, err = scanDuration([]byte("3000000000"))
	require.NoError(t, err)
	assert.Equal(t, 3*time.Second, value)

	_, err = scanDuration(struct{}{})
	require.Error(t, err)
	_, err = scanDuration("not-a-duration")
	require.Error(t, err)
}

func TestScalarConversions(t *testing.T) {
	stringValue, ok := scalarString(reflect.ValueOf(true))
	require.True(t, ok)
	assert.Equal(t, "true", stringValue)

	stringValue, ok = scalarString(reflect.ValueOf(float32(1.5)))
	require.True(t, ok)
	assert.Equal(t, "1.5", stringValue)

	_, ok = scalarString(reflect.ValueOf(struct{}{}))
	assert.False(t, ok)

	boolValue, ok := scalarBool(reflect.ValueOf(int64(1)))
	require.True(t, ok)
	assert.True(t, boolValue)

	boolValue, ok = scalarBool(reflect.ValueOf(uint64(0)))
	require.True(t, ok)
	assert.False(t, boolValue)

	boolValue, ok = scalarBool(reflect.ValueOf(float64(0.25)))
	require.True(t, ok)
	assert.True(t, boolValue)

	_, ok = scalarBool(reflect.ValueOf("true"))
	assert.False(t, ok)

	intValue, ok := scalarInt(reflect.ValueOf(uint64(10)))
	require.True(t, ok)
	assert.EqualValues(t, 10, intValue)

	_, ok = scalarInt(reflect.ValueOf(uint64(1 << 63)))
	assert.False(t, ok)

	intValue, ok = scalarInt(reflect.ValueOf(1.9))
	require.True(t, ok)
	assert.EqualValues(t, 1, intValue)

	uintValue, ok := scalarUint(reflect.ValueOf(int64(10)))
	require.True(t, ok)
	assert.EqualValues(t, 10, uintValue)

	_, ok = scalarUint(reflect.ValueOf(int64(-1)))
	assert.False(t, ok)

	uintValue, ok = scalarUint(reflect.ValueOf(2.9))
	require.True(t, ok)
	assert.EqualValues(t, 2, uintValue)

	_, ok = scalarUint(reflect.ValueOf(-2.9))
	assert.False(t, ok)

	floatValue, ok := scalarFloat(reflect.ValueOf(uint64(10)))
	require.True(t, ok)
	assert.Equal(t, float64(10), floatValue)

	_, ok = scalarFloat(reflect.ValueOf("1.2"))
	assert.False(t, ok)
}

func TestScanPostgresArrayElement(t *testing.T) {
	value, err := scanPostgresArrayElement(postgresArrayElement{value: "true"}, reflect.TypeOf(false))
	require.NoError(t, err)
	assert.True(t, value.Bool())

	value, err = scanPostgresArrayElement(postgresArrayElement{value: "42"}, reflect.TypeOf(uint16(0)))
	require.NoError(t, err)
	assert.EqualValues(t, 42, value.Uint())

	value, err = scanPostgresArrayElement(postgresArrayElement{value: "1.5"}, reflect.TypeOf(float32(0)))
	require.NoError(t, err)
	assert.InDelta(t, 1.5, value.Float(), 0.0001)

	value, err = scanPostgresArrayElement(postgresArrayElement{value: "2024-01-02T03:04:05Z"}, reflect.TypeOf(time.Time{}))
	require.NoError(t, err)
	assert.Equal(t, 2024, value.Interface().(time.Time).Year())

	value, err = scanPostgresArrayElement(postgresArrayElement{value: "value"}, reflect.TypeOf((*string)(nil)))
	require.NoError(t, err)
	assert.Equal(t, "value", *value.Interface().(*string))

	value, err = scanPostgresArrayElement(postgresArrayElement{null: true}, reflect.TypeOf((*string)(nil)))
	require.NoError(t, err)
	assert.True(t, value.IsNil())

	type payload struct {
		Name string `json:"name"`
	}
	value, err = scanPostgresArrayElement(postgresArrayElement{value: `{"name":"test"}`}, reflect.TypeOf(payload{}))
	require.NoError(t, err)
	assert.Equal(t, "test", value.Interface().(payload).Name)
}

func TestScanSQLValueErrorsAndValuer(t *testing.T) {
	_, err := scanSQLValue[int](failingValuer{})
	require.Error(t, err)

	_, err = scanSQLValue[int](struct{}{})
	require.Error(t, err)

	_, _, err = scanPostgresArrayText[[2]int]("{1}", reflect.TypeOf([2]int{}))
	require.Error(t, err)

	_, _, err = scanPostgresArrayText[[]int](`{"unterminated}`, reflect.TypeOf([]int{}))
	require.Error(t, err)

	_, err = scanText[complex64]("not-complex", reflect.TypeOf(complex64(0)), []byte("not-complex"))
	require.Error(t, err)
}
