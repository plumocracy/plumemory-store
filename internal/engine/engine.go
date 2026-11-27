// package engine provides the storage engine.
package engine

import (
	"errors"
	"reflect"
	"sync"
)

// Just to start, there are at least a few things I need to investigate here.
// 1: Are unlocked reads safe? Should reads also be locked? I think that reads can be safely done unlocked
// except for in instances where a read and write of the same value occur at the same time across different
// go routines. Is the write lock enough?
// 2: Is data[any]any efficient enough? Should I be doing more manual memory management?
// 3: I need to add a delete method.

type Engine struct {
	mtx  sync.Mutex
	data map[any]any
}

func New() *Engine {
	return &Engine{
		data: make(map[any]any),
	}
}

// Stores a key value pair in memory.
// All key types are valid except func, chan. Key must not be nil.
func (e *Engine) Store(key any, value any) error {
	err := isInterfaceAnyValidType(key)
	if err != nil {
		return err
	}

	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.data[key] = value

	return nil
}

func (e *Engine) Get(key any) (any, error) {
	err := isInterfaceAnyValidType(key)
	if err != nil {
		return nil, err
	}

	v, ok := e.data[key]
	if !ok {
		return nil, errors.New("Provided key has no stored value.")
	}

	return v, nil
}

func isInterfaceAnyValidType(i interface{}) error {
	if i == nil {
		return errors.New("Key must not be nil.")
	}

	v := reflect.ValueOf(i)

	switch v.Kind() {
	case reflect.Chan:
		return errors.New("Type chan cannot be a key.")
	case reflect.Func:
		return errors.New("Type func cannot be a key.")
	case reflect.Ptr, reflect.Map, reflect.Slice,
		reflect.Interface, reflect.UnsafePointer:
		if !v.IsNil() {
			return nil
		} else {
			return errors.New("Key must not be nil.")
		}
	}

	return nil
}
