package assert

import (
	"errors"
	"fmt"
	"io"
)

// Int interface type
type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Uint interface type
type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Xint interface type. alias of Integer
type Xint interface {
	Int | Uint
}

// Integer interface type. all int or uint types
type Integer interface {
	Int | Uint
}

// Float interface type
type Float interface {
	~float32 | ~float64
}

// Number interface type. contains all int, uint and float types
type Number interface {
	Int | Uint | Float
}

// ScalarType basic interface type.
//
// TIP: has bool type, it cannot be ordered
//
// contains: (x)int, float, ~string, ~bool types
type ScalarType interface {
	Int | Uint | Float | ~string | ~bool
}

// ErrConvType error
var ErrConvType = errors.New("convert value type error")

// Int64able interface
type Int64able interface {
	Int64() (int64, error)
}

// Float64able interface
type Float64able interface {
	Float64() (float64, error)
}

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Helper()
	Name() string
	Error(args ...any)
}

type failNower interface {
	FailNow()
}

// StringWriteStringer interface
type StringWriteStringer interface {
	io.StringWriter
	fmt.Stringer
}
