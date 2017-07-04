package eff

import "runtime"
import "fmt"

const (
	Field_Msg   string = "msg"
	Field_Error string = "error"
	Field_Stack string = "stack"
)

// capturing a stack trace requires allocating a slice of bytes before we know how much space we'll actually need
const stackMaxSizeBytes = 1 << 14
const stackIncludesAllGoRoutines = false

// Effer is the interface used by Eff
type Effer interface {
	// Error is the std go error interface
	Error() string

	// ExitCode can be used as the process exit code, for when this error should end the process (works with urfave/cli)
	ExitCode() int

	// Fields gets a copy of the underlying field map
	Fields() map[string]interface{}

	// Stack gets the (formatted) stack trace from this Effer
	Stack() string

	// WithField returns a new Effer with the specified field/value added
	WithField(name string, value interface{}) Effer

	// WithFields returns a new Effer with the given fields added
	WithFields(fields map[string]interface{}) Effer

	// WithMessage returns a new Effer using the given message
	WithMessage(msg string) Effer

	// WithExitCode returns a new Effer that uses this exit code
	WithExitCode(code int) Effer
}

// Eff -- Errors with Fields Forever
//
// Err is an attempt to allow anyone to add context to an error that is being passed up a chain of calls, so that by
// the time someone actually handles that error, the trip it took to get there is not forgotten.
//
//  - Fields are specific key/value pairs, and setting the same key twice destroys the original
//  - Stack is the stack trace at the moment eff.New() was called.
//
// Note that the With* funcs do not alter the calling object, but instead copy everything into a new instance.
// For example, if 'e' has field 'test':'value', and you call f := e.WithField('date', 1234), the new 'f' has
// the date field, but the original object 'e' does not.
//
type Eff struct {
	fields   map[string]interface{}
	exitCode int
}

// New returns a new Eff object, ready to use
func New() Effer {
	e := &Eff{
		fields:   make(map[string]interface{}),
		exitCode: -1,
	}
	buf := make([]byte, stackMaxSizeBytes)
	bufLen := runtime.Stack(buf, stackIncludesAllGoRoutines)
	stack := string(buf[0 : bufLen-1])
	if bufLen >= stackMaxSizeBytes {
		stack += "|STACK CUT OFF, try increasing stackMaxSizeBytes"
	}
	e.fields[Field_Stack] = stack
	return e
}

func (e *Eff) copy() *Eff {
	f := &Eff{
		fields:   make(map[string]interface{}),
		exitCode: e.exitCode,
	}
	for k, v := range e.fields {
		f.fields[k] = v
	}
	return f
}

// NewErr allows you to wrap an existing error.
// If the error is a *Eff, it will be copied into the new *Eff.
// If the error is some other type, it will be set to Field_Error.
func NewErr(err error) Effer {
	switch v := err.(type) {
	case *Eff:
		return v.copy()
	default:
		return New().WithField(Field_Error, err)
	}
}

// NewMsg is simply a shorter version of New().WithMessage(msg)
func NewMsg(msg string) Effer {
	return New().WithMessage(msg)
}

// Error is an implementation of golang's error interface
func (e *Eff) Error() string {
	// TODO: build string including fields?
	return fmt.Sprintf("%v", e.Fields())
}

// ExitCode specifies the process exit code that should be used if this error ends the process.
// Defaults to -1.
func (e *Eff) ExitCode() int {
	return e.exitCode
}

// Fields returns a copy of all fields (not including stack trace or exit code)
func (e *Eff) Fields() map[string]interface{} {
	fields := make(map[string]interface{})
	for k, v := range e.fields {
		if k == Field_Stack {
			continue
		}
		fields[k] = v
	}
	return fields
}

func (e *Eff) Stack() string {
	return e.fields[Field_Stack].(string)
}

// WithField adds a single field to itself
func (e *Eff) WithField(name string, value interface{}) Effer {
	f := e.copy()
	f.fields[name] = value
	return f
}

// WithFields adds all passed in fields to itself
func (e *Eff) WithFields(fields map[string]interface{}) Effer {
	f := e.copy()
	for k, v := range fields {
		f.fields[k] = v
	}
	return f
}

func (e *Eff) WithMessage(msg string) Effer {
	return e.WithField(Field_Msg, msg)
}

func (e *Eff) WithExitCode(code int) Effer {
	f := e.copy()
	f.exitCode = code
	return f
}
