package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/plumocracy/inmem/internal/engine"
)

// This is an ugly and kind of bad repl, but this isn't about the repl.

// Repl functions
// store(key, value)
// read(key)
// delete(key)
// exit()

type Repl struct {
	Reader   io.Reader
	Writer   io.Writer
	Commands map[string]func([]string, *engine.Engine) error
}

type ParsedCommand struct {
	Name string
	Args []string
}

func New() *Repl {
	return &Repl{
		Reader:   os.Stdin,
		Commands: make(map[string]func([]string, *engine.Engine) error),
	}
}

func (r *Repl) RegisterCommand(name string, callback func([]string, *engine.Engine) error) error {
	r.Commands[name] = callback
	return nil
}

func (r *Repl) Run(e *engine.Engine) error {
	rdr := bufio.NewReader(r.Reader)

	for {
		fmt.Print("> ")

		str, err := rdr.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		pc, err := ParseLine(str)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		callback, ok := r.Commands[pc.Name]
		if !ok {
			fmt.Println("Invalid Command")
			continue
		}

		// These callbacks will be responsble for printing to console
		// and everything. A better arch may be returning some value (interface{}?)
		// and printing it here.
		err = callback(pc.Args, e)
		if err != nil {
			// Yeah this is fucking ugly idc rn
			switch err.Error() {
			case "exit":
				os.Exit(0)
			default:
				fmt.Println(err)
			}
		}
	}
}

// A line is always in format <command>(<options>)
func ParseLine(line string) (*ParsedCommand, error) {
	argsStart := -1
	argsEnd := -1

	for i, c := range line {
		switch c {
		case '(':
			argsStart = i + 1
		case ')':
			argsEnd = i
		}
	}

	// Special case "exit"
	if strings.Contains(line, "exit") {
		os.Exit(0)
	}

	if argsStart == -1 || argsEnd == -1 {
		return nil, errors.New("Malformed command")
	}

	name := line[:argsStart-1]
	args := line[argsStart:argsEnd]

	// Note that some args may have a leading space.
	// Consumers are expected to deal with this.
	argsSlice := strings.Split(args, ",")

	pc := &ParsedCommand{
		Name: name,
		Args: argsSlice,
	}

	return pc, nil
}
