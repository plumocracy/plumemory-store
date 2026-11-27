package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/plumocracy/inmem/internal/engine"
	"github.com/plumocracy/inmem/internal/repl"
)

func main() {
	r := repl.New()
	e := engine.New()

	r.RegisterCommand("store", StoreValue)
	r.RegisterCommand("read", GetValue)
	r.RegisterCommand("delete", DeleteValue)
	r.RegisterCommand("exit", CloseRepl)

	err := r.Run(e)
	if err != nil {
		fmt.Println(err)
	}
}

func StoreValue(args []string, e *engine.Engine) error {
	if len(args) < 2 {
		return errors.New("Not enough arguments to store.")
	} else if len(args) > 2 {
		return errors.New("To many arguments to store.")
	}

	for _, str := range args {
		strings.Trim(str, " ")
	}

	err := e.Store(args[0], args[1])
	if err != nil {
		return err
	}

	fmt.Printf("Stored %v as %v\n", args[0], args[1])
	return nil
}

func GetValue(args []string, e *engine.Engine) error {
	if len(args) != 1 {
		return errors.New("Malformed argument count.")
	}

	value, err := e.Read(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Got value %v\n", value)
	return nil
}

func DeleteValue(args []string, e *engine.Engine) error {
	if len(args) != 1 {
		return errors.New("Malformed argument count.")
	}

	err := e.Delete(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Deleted value at %v\n", args[0])
	return nil
}

func CloseRepl(args []string, e *engine.Engine) error {
	return errors.New("exit")
}
