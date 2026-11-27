package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInterfaceAnyValidType(t *testing.T) {
	var err error

	err = isInterfaceAnyValidType(nil)
	assert.EqualError(t, err, NilKeyMsg)
	err = nil

	ch := make(chan string)
	err = isInterfaceAnyValidType(ch)
	assert.EqualError(t, err, ChanKeyMsg)

	var p *int
	err = isInterfaceAnyValidType(p)
	assert.EqualError(t, err, NilKeyMsg)

	err = nil
	err = isInterfaceAnyValidType(func() {})
	assert.EqualError(t, err, FuncKeyMsg)

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
	name, err := e.Read("name")
	assert.NoError(t, err)
	assert.Equal(t, "plum", name)

	name, err = e.Read("fart")
	assert.EqualError(t, err, MissingKeyMsg)
}

func TestDeleteData(t *testing.T) {
	e := New()
	e.Store("name", "plum")
	name, err := e.Read("name")
	assert.NoError(t, err)
	assert.Equal(t, "plum", name)

	err = e.Delete("name")
	assert.NoError(t, err)

	_, err = e.Read("name")
	assert.EqualError(t, err, MissingKeyMsg)
}
