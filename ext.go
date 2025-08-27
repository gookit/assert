package assert

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
)

//
// region Buffer
//

// Buffer for unit testing
type Buffer struct {
	bytes.Buffer
}

// NewBuffer instance
func NewBuffer() *Buffer { return &Buffer{} }

// ResetGet get buffer string and reset.
func (b *Buffer) ResetGet() string {
	s := b.String()
	b.Reset()
	return s
}

//
// region Thread-safe buffer
//

// SafeBuffer Thread-safe buffer for testing
type SafeBuffer struct {
	bytes.Buffer
	mu sync.Mutex
}

// NewSafeBuffer create a new SafeBuffer
func NewSafeBuffer() *SafeBuffer {
	return &SafeBuffer{}
}

// Write implements io.Writer
func (sb *SafeBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.Buffer.Write(p)
}

// WriteByte override parent WriteByte method
func (sb *SafeBuffer) WriteByte(b byte)  error {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.Buffer.WriteByte(b)
}

// WriteRune override parent WriteRune method
func (sb *SafeBuffer) WriteRune(r rune)  (int,error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.Buffer.WriteRune(r)
}

// WriteString implements io.StringWriter
func (sb *SafeBuffer) WriteString(s string) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.Buffer.WriteString(s)
}

// WriteTo override parent WriteTo method
func (sb *SafeBuffer) WriteTo(w io.Writer) (n int64, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.Buffer.WriteTo(w)
}

// ResetGet get buffer content and reset
func (sb *SafeBuffer) ResetGet() string {
	s := sb.String()
	sb.Reset()
	return s
}

//
// region Mock environment
//

// backup os ENV
var envBak = os.Environ()

// MockEnvValue will store old env value, set new val. will restore old value on end.
func MockEnvValue(key, val string, fn func(nv string)) {
	old := os.Getenv(key)
	err := os.Setenv(key, val)
	if err != nil {
		panic(err)
	}

	// call with new value
	fn(os.Getenv(key))

	// if old is empty, unset key.
	if old == "" {
		err = os.Unsetenv(key)
	} else {
		err = os.Setenv(key, old)
	}
	if err != nil {
		panic(err)
	}
}

// MockOsEnv mock clean os.Environ by set input map
//
//  - will CLEAR all old ENV data, use given a data map.
//  - will RECOVER old ENV after fn run.
func MockOsEnv(mp map[string]string, fn func()) {
	os.Clearenv()
	for key, val := range mp {
		_ = os.Setenv(key, val)
	}

	fn()

	os.Clearenv()
	for _, str := range envBak {
		nodes := strings.SplitN(str, "=", 2)
		_ = os.Setenv(nodes[0], nodes[1])
	}
}

// MockOsEnvByText mock clean os Environ by input text lines. see MockOsEnv
//
// Usage:
//
//	assert.MockOsEnvByText(`
//		APP_COMMAND = login
//		APP_ENV = dev
//		APP_DEBUG = true
//	`, func() {
//			// do something ...
//	})
func MockOsEnvByText(envText string, fn func()) {
	ss := strings.Split(envText, "\n")
	mp := make(map[string]string, len(ss))

	for _, line := range ss {
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		if line[0] == '#' || strings.HasPrefix(line, "//") {
			continue
		}

		nodes := strings.SplitN(line, "=", 2)
		envKey := strings.TrimSpace(nodes[0])

		if len(nodes) < 2 {
			mp[envKey] = ""
		} else {
			mp[envKey] = strings.TrimSpace(nodes[1])
		}
	}

	MockOsEnv(mp, fn)
}