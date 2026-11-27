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

const (
	MissingKeyMsg = "provided key has no value"
	NilKeyMsg     = "key must not be nil"
	FuncKeyMsg    = "kust must not be of type func"
	ChanKeyMsg    = "kust must not be of type chan"
)

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

func (e *Engine) Read(key any) (any, error) {
	err := isInterfaceAnyValidType(key)
	if err != nil {
		return nil, err
	}

	v, ok := e.data[key]
	if !ok {
		return nil, errors.New(MissingKeyMsg)
	}

	return v, nil
}

func (e *Engine) Delete(key any) error {
	err := isInterfaceAnyValidType(key)
	if err != nil {
		return err
	}

	e.mtx.Lock()
	defer e.mtx.Unlock()
	_, ok := e.data[key]
	if ok {
		delete(e.data, key)
		return nil
	} else {
		return errors.New(NilKeyMsg)
	}
}

func isInterfaceAnyValidType(i interface{}) error {
	if i == nil {
		return errors.New(NilKeyMsg)
	}

	v := reflect.ValueOf(i)

	switch v.Kind() {
	case reflect.Chan:
		return errors.New(ChanKeyMsg)
	case reflect.Func:
		return errors.New(FuncKeyMsg)
	case reflect.Ptr, reflect.Map, reflect.Slice,
		reflect.Interface, reflect.UnsafePointer:
		if !v.IsNil() {
			return nil
		} else {
			return errors.New(NilKeyMsg)
		}
	}

	return nil
}
