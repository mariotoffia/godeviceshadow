package jsonutils_test

import (
	"reflect"
	"testing"

	"github.com/mariotoffia/godeviceshadow/utils/jsonutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Address struct {
	City  string `json:"city"`
	State string `json:"state"`
	Zip   int    `json:"zip"`
}
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestWhenErrorInMiddleOfPayload(t *testing.T) {
	data := []byte(`{"name": "john", "age": 30, "address": {"city": "new york" "state": "ny", "zip": 100}`)
	_, err := jsonutils.UnmarshalJSON(data, reflect.TypeOf(Person{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), `^ invalid character '"' after object key:value pair`)
}
