package parse

import (
	"fmt"
	"sync"
)

type ErrorKind int

const (
	Error ErrorKind = iota
	Warning
	Note
)

func (k ErrorKind) String() string {
	switch k {
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Note:
		return "note"
	}
	return ""
}

type Err struct {
	File string
	Len  int
	Loc  TokenLocation
	Msg  string
	Kind ErrorKind
}

func (e Err) String() string {
	startLine := e.Loc.Start.Line + 1
	startCol := e.Loc.Start.Col + 1
	fileName := e.File

	return fmt.Sprintf("%s:%d:%d: %s: %s", fileName, startLine, startCol, e.Kind, e.Msg)
}

type ErrorContainer struct {
	Errors []Err
	mu     *sync.Mutex
}

func (e *ErrorContainer) AddError(err Err) {
	e.mu.Lock()
	e.Errors = append(e.Errors, err)
	e.mu.Unlock()
}

func (e *ErrorContainer) GetErrors() []Err {
	return e.Errors
}

func (e *ErrorContainer) Clear() {
	e.Errors = []Err{}
}

func (e *ErrorContainer) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e *ErrorContainer) HasWarnings() bool {
	for _, err := range e.Errors {
		if err.Kind == Warning {
			return true
		}
	}
	return false
}

func (e *ErrorContainer) HasNotes() bool {
	for _, err := range e.Errors {
		if err.Kind == Note {
			return true
		}
	}
	return false
}

func (e *ErrorContainer) Print() {
	for _, err := range e.Errors {
		println(err.String())
	}
}
