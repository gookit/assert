package assert_test

import (
	"os"
	"testing"

	"github.com/gookit/assert"
)

func TestBuffer(t *testing.T) {
	// SafeBuffer
	sb := assert.NewSafeBuffer()
	_, err := sb.Write([]byte("hello,"))
	assert.NoErr(t, err)
	err = sb.WriteByte('a')
	assert.NoErr(t, err)
	_, err = sb.WriteRune('b')
	assert.NoErr(t, err)

	s := sb.ResetGet()
	assert.Eq(t, "hello,ab", s)

	_, err = sb.WriteString("hello")
	assert.NoErr(t, err)

	// Buffer
	b := assert.NewBuffer()
	_, err = sb.WriteTo(b)
	assert.NoErr(t, err)
	assert.Eq(t, "hello", b.ResetGet())
}

func TestMockEnvValue(t *testing.T) {
	is := assert.New(t)
	is.Eq("", os.Getenv("APP_COMMAND"))

	assert.MockEnvValue("APP_COMMAND", "new val", func(nv string) {
		is.Eq("new val", nv)
	})

	shellVal := "custom-value"
	assert.MockEnvValue("SHELL", shellVal, func(newVal string) {
		is.Eq(shellVal, newVal)
	})

	is.Eq("", os.Getenv("APP_COMMAND"))
	is.Panics(func() {
		assert.MockEnvValue("invalid=", "value", nil)
	})
}

func TestMockOsEnv(t *testing.T) {
	is := assert.New(t)
	is.Eq("", os.Getenv("APP_COMMAND"))

	assert.MockOsEnv(map[string]string{
		"APP_COMMAND": "new val",
	}, func() {
		is.Eq("new val", os.Getenv("APP_COMMAND"))
	})

	is.Eq("", os.Getenv("APP_COMMAND"))
}

func TestMockOsEnvByText(t *testing.T) {
	envStr := `
APP_COMMAND = login
APP_ENV = dev
APP_DEBUG = true
APP_PWD=
// ENV_NOT_EXIST=comment line
`

	assert.MockOsEnvByText(envStr, func() {
		assert.Len(t, os.Environ(), 4)
		assert.Eq(t, "true", os.Getenv("APP_DEBUG"))
		assert.Eq(t, "login", os.Getenv("APP_COMMAND"))
		assert.Empty(t, os.Getenv("APP_PWD"))
		assert.Empty(t, os.Getenv("ENV_NOT_EXIST"))
	})
}