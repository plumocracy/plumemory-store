package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInterfaceAnyValidType(t *testing.T) {
	var err error

	err = isInterfaceAnyValidType(nil)
	assert.Error(t, err)
	err = nil

	ch := make(chan string)
	err = isInterfaceAnyValidType(ch)
	assert.Error(t, err)

	var p *int
	err = isInterfaceAnyValidType(p)
	assert.Error(t, err)

	err = nil
	err = isInterfaceAnyValidType(func() {})
	assert.Error(t, err)

	err = isInterfaceAnyValidType("test")
	assert.NoError(t, err)

	err = isInterfaceAnyValidType(1234)
	assert.NoError(t, err)

	err = isInterfaceAnyValidType(13.5)
	assert.NoError(t, err)

}

func TestStoreAndRetrieveData(t *testing.T) {
	e := New()
	e.Store("name", "plum")
	name, err := e.Get("name")
	assert.NoError(t, err)
	assert.Equal(t, "plum", name)

	name, err = e.Get("fart")
	assert.Error(t, err)
}
