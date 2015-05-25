package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/zubairhamed/go-lwm2m/core"
	"testing"
	"time"
	"github.com/zubairhamed/go-lwm2m"
	"github.com/zubairhamed/go-lwm2m/registry"
	"github.com/zubairhamed/go-lwm2m/objects/oma"
)

func TestGetValueByteLength(t *testing.T) {
	test1 := [] struct {
		input 		interface{}
		expected	uint32
	}{
		{-128, 1},
		{127, 1},
		{-32768, 2},
		{-2147483648, 4},
		{2147483647, 4},
		{-9223372036854775808, 8},
		{9223372036854775807, 8},
		{-3.4E+38, 4},
		{3.4E+38, 4},
		{-1.7E+308, 8},
		{1.7E+308, 8},
		{"this is a string", 16},
		{true, 1},
		{false, 1},
		{time.Now(), 8},
		{[]byte{}, 0},
	}

	for _, c := range test1 {
		v, _ := core.GetValueByteLength(c.input)
		assert.Equal(t, v, uint32(c.expected), "Wrong expected length returned")
	}

	test2 := [] struct {
		input 		interface{}
	}{
		{uint(1)},
		{uint16(1)},
		{uint32(1)},
		{uint64(1)},
	}

	for _, c := range test2 {
		_, err := core.GetValueByteLength(c.input)
		assert.NotNil(t, err, "An error should be returned")
	}
}

func TestObjectData(t *testing.T) {

	tests := [] struct {
		path 	string
		value 	interface{}
	}{
		{"/0/0", 1},
		{"/0/1", 0},
		{"/0/2/101", []byte{0, 15}},
		{"/0/3", 101},
		{"/1/0", 1},
		{"/1/1", 1},
		{"/1/2/102", []byte{0, 15}},
		{"/1/3", 102},
		{"/2/0", 3},
		{"/2/1", 0},
		{"/2/2/101", []byte{0, 15}},
		{"/2/2/102", []byte{0, 1}},
		{"/2/3", 101},
		{"/3/0", 4},
		{"/3/1", 0},
		{"/3/2/101", []byte{0, 1}},
		{"/3/2/0", []byte{0, 1}},
		{"/3/3", 101},
		{"/4/0", 5},
		{"/4/1", 65535},
		{"/4/2/101", []byte{0, 16}},
		{"/4/3", 65535},
	}

	data := &core.ObjectsData{
		Data: make(map[string]interface{}),
	}

	for _, c := range tests {
		data.Put(c.path, c.value)
		assert.Equal(t, data.Get(c.path), c.value, "Value get not equal to put: (",  c.path, "vs", c.value)
	}

	assert.Equal(t, data.Length(), 22, "Number of items in ObjectData. Expected", 22, "actual", data.Length())

	data.Clear()
	assert.Equal(t, data.Length(), 0, "Number of items in ObjectData. Expected", 0, "actual", data.Length())
}

func TestBuildResourceStringPayload(t *testing.T) {
	client := lwm2m.NewLWM2MClient(":0", "localhost:5683")

	reg := registry.NewDefaultObjectRegistry()
	client.UseRegistry(reg)

	client.EnableObject(oma.OBJECT_LWM2M_SECURITY, nil)
	client.EnableObject(oma.OBJECT_LWM2M_ACCESS_CONTROL, nil)
	client.EnableObject(oma.OBJECT_LWM2M_CONNECTIVITY_MONITORING, nil)

	str := core.BuildModelResourceStringPayload(client.GetEnabledObjects())
	compare := "</0>,</2>,</4>,"

	assert.Equal(t, str, compare, "Unexpected output building Model Resource String")
}