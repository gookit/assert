# Assert

✅ 工具包 `gookit/assert` 提供了一些常用的 Go 单元测试中用于断言的工具函数。

- 使用非常简单 eg: `assert.Eq(t, want, give)` `assert.Err(t, err)`
- 非常轻量，没有外部依赖库
- 支持 go `1.18+`

> 本身代码来自于 [gookit/goutil](https://github.com/gookit/goutil) 下的 `testutil/assert` 包

## Install

```bash
go get github.com/gookit/assert
```

## GoDocs

Please see [Go docs](https://pkg.go.dev/github.com/gookit/assert)

## Usage

```go
package assert_test

import (
	"testing"

	"github.com/gookit/assert"
	"github.com/gookit/goutil/errorx"
)

func TestErr(t *testing.T) {
	err := errorx.Raw("this is a error")

	assert.NoErr(t, err, "user custom message")
	assert.ErrMsg(t, err, "this is a error")
}
```

## Function API

> generate by: `go doc .`

```go
func Contains(t TestingT, src, elem any, fmtAndArgs ...any) bool
func ContainsElems[T ScalarType](t TestingT, list, sub []T, fmtAndArgs ...any) bool
func ContainsKey(t TestingT, mp, key any, fmtAndArgs ...any) bool
func ContainsKeys(t TestingT, mp any, keys any, fmtAndArgs ...any) bool
func DirExists(t TestingT, dirPath string, fmtAndArgs ...any) bool
func DirNotExists(t TestingT, dirPath string, fmtAndArgs ...any) bool
func Empty(t TestingT, give any, fmtAndArgs ...any) bool
func Eq(t TestingT, want, give any, fmtAndArgs ...any) bool
func Equal(t TestingT, want, give any, fmtAndArgs ...any) bool
func Err(t TestingT, err error, fmtAndArgs ...any) bool
func ErrIs(t TestingT, err, wantErr error, fmtAndArgs ...any) bool
func ErrMsg(t TestingT, err error, wantMsg string, fmtAndArgs ...any) bool
func ErrMsgContains(t TestingT, err error, subMsg string, fmtAndArgs ...any) bool
func ErrSubMsg(t TestingT, err error, subMsg string, fmtAndArgs ...any) bool
func Error(t TestingT, err error, fmtAndArgs ...any) bool
func Fail(t TestingT, failMsg string, fmtAndArgs ...any) bool
func FailNow(t TestingT, failMsg string, fmtAndArgs ...any) bool
func False(t TestingT, give bool, fmtAndArgs ...any) bool
func FileExists(t TestingT, filePath string, fmtAndArgs ...any) bool
func FileNotExists(t TestingT, filePath string, fmtAndArgs ...any) bool
func Gt(t TestingT, give, min any, fmtAndArgs ...any) bool
func Gte(t TestingT, give, min any, fmtAndArgs ...any) bool
func HideFullPath()
func IsKind(t TestingT, wantKind reflect.Kind, give any, fmtAndArgs ...any) bool
func IsType(t TestingT, wantType, give any, fmtAndArgs ...any) bool
func Len(t TestingT, give any, wantLn int, fmtAndArgs ...any) bool
func LenGt(t TestingT, give any, minLn int, fmtAndArgs ...any) bool
func Lt(t TestingT, give, max any, fmtAndArgs ...any) bool
func Lte(t TestingT, give, max any, fmtAndArgs ...any) bool
func Neq(t TestingT, want, give any, fmtAndArgs ...any) bool
func Nil(t TestingT, give any, fmtAndArgs ...any) bool
func NoErr(t TestingT, err error, fmtAndArgs ...any) bool
func NoError(t TestingT, err error, fmtAndArgs ...any) bool
func NotContains(t TestingT, src, elem any, fmtAndArgs ...any) bool
func NotContainsKey(t TestingT, mp, key any, fmtAndArgs ...any) bool
func NotContainsKeys(t TestingT, mp any, keys any, fmtAndArgs ...any) bool
func NotEmpty(t TestingT, give any, fmtAndArgs ...any) bool
func NotEq(t TestingT, want, give any, fmtAndArgs ...any) bool
func NotEqual(t TestingT, want, give any, fmtAndArgs ...any) bool
func NotNil(t TestingT, give any, fmtAndArgs ...any) bool
func NotPanics(t TestingT, fn PanicRunFunc, fmtAndArgs ...any) bool
func NotSame(t TestingT, want, actual any, fmtAndArgs ...any) bool
func Panics(t TestingT, fn PanicRunFunc, fmtAndArgs ...any) bool
func PanicsErrMsg(t TestingT, fn PanicRunFunc, errMsg string, fmtAndArgs ...any) bool
func PanicsMsg(t TestingT, fn PanicRunFunc, wantVal any, fmtAndArgs ...any) bool
func Same(t TestingT, wanted, actual any, fmtAndArgs ...any) bool
func StrContains(t TestingT, s, sub string, fmtAndArgs ...any) bool
func StrContainsAll(t TestingT, s string, subs []string, fmtAndArgs ...any) bool
func StrCount(t TestingT, s, sub string, count int, fmtAndArgs ...any) bool
func StrNotContains(t TestingT, s, sub string, fmtAndArgs ...any) bool
func True(t TestingT, give bool, fmtAndArgs ...any) bool
```

## Code Check & Testing

```bash
gofmt -w -l ./
# testing
go test ./...
```

## Refer

- https://github.com/gookit/goutil/tree/master/testutil/assert
- https://github.com/stretchr/testify