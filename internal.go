package assert

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// isNil value check
func checkIsNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

// isEmpty value check
func isEmpty(v any) bool {
	if v == nil {
		return true
	}
	return reflectIsEmpty(reflect.ValueOf(v))
}

func checkContains(data, elem any) (valid, found bool) {
	if data == nil {
		return false, false
	}

	dataRv := reflect.ValueOf(data)
	dataRt := reflect.TypeOf(data)
	dataKind := dataRt.Kind()

	// string
	if dataKind == reflect.String {
		return true, strings.Contains(dataRv.String(), fmt.Sprint(elem))
	}

	// map
	if dataKind == reflect.Map {
		mapKeys := dataRv.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if reflectIsEqual(mapKeys[i].Interface(), elem) {
				return true, true
			}
		}
		return true, false
	}

	// array, slice - other return false
	if dataKind != reflect.Slice && dataKind != reflect.Array {
		return false, false
	}

	for i := 0; i < dataRv.Len(); i++ {
		if reflectIsEqual(dataRv.Index(i).Interface(), elem) {
			return true, true
		}
	}
	return true, false
}

//
// region math utils
// - from github.com/gookit/goutil/mathutil
//

// MaxInt compare and return max value
func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Compare any intX,floatX value by given op. returns `first op(=,!=,<,<=,>,>=) second`
//
// Usage:
//
//	mathutil.Compare(2, 3, ">") // false
//	mathutil.Compare(2, 1.3, ">") // true
//	mathutil.Compare(2.2, 1.3, ">") // true
//	mathutil.Compare(2.1, 2, ">") // true
func mathCompare(first, second any, op string) bool {
	if first == nil || second == nil {
		return false
	}

	switch fVal := first.(type) {
	case float64:
		if sVal, err := toFloat(second); err == nil {
			return mathCompFloat(fVal, sVal, op)
		}
	case float32:
		if sVal, err := toFloat(second); err == nil {
			return mathCompFloat(float64(fVal), sVal, op)
		}
	default: // as int64
		if int1, err := toInt64(first); err == nil {
			if int2, err2 := toInt64(second); err2 == nil {
				return mathCompInt64(int1, int2, op)
			}
		}
	}

	return false
}

// CompInt64 compare int64 value. returns `first op(=,!=,<,<=,>,>=) second`
func mathCompInt64(first, second int64, op string) bool {
	return mathCompValue(first, second, op)
}

// CompFloat compare float64,float32 value. returns `first op(=,!=,<,<=,>,>=) second`
func mathCompFloat[T Float](first, second T, op string) (ok bool) {
	return mathCompValue(first, second, op)
}

// CompValue compare intX,uintX,floatX value. returns `first op(=,!=,<,<=,>,>=) second`
func mathCompValue[T Number](first, second T, op string) (ok bool) {
	switch op {
	case "<", "lt":
		ok = first < second
	case "<=", "lte":
		ok = first <= second
	case ">", "gt":
		ok = first > second
	case ">=", "gte":
		ok = first >= second
	case "=", "eq":
		ok = first == second
	case "!=", "ne", "neq":
		ok = first != second
	}
	return
}

// ToInt64 convert value to int64, return error on failed
func toInt64(in any) (int64, error) { return toInt64With(in) }

// ToInt64With try to convert value to int64. can with some option func, more see ConvOption.
func toInt64With(in any) (i64 int64, err error) {
	if in == nil {
		return 0, nil
	}

	switch tVal := in.(type) {
	case string:
		sVal := strings.TrimSpace(tVal)
		i64, err = strconv.ParseInt(sVal, 10, 0)
		// handle the case where the string might be a float
		if err != nil && strIsNumeric(sVal) {
			var floatVal float64
			if floatVal, err = strconv.ParseFloat(sVal, 64); err == nil {
				i64 = int64(math.Round(floatVal))
				err = nil
			}
		}
	case int:
		i64 = int64(tVal)
	case int8:
		i64 = int64(tVal)
	case int16:
		i64 = int64(tVal)
	case int32:
		i64 = int64(tVal)
	case int64:
		i64 = tVal
	case *int64: // default support int64 ptr type
		i64 = *tVal
	case uint:
		i64 = int64(tVal)
	case uint8:
		i64 = int64(tVal)
	case uint16:
		i64 = int64(tVal)
	case uint32:
		i64 = int64(tVal)
	case uint64:
		i64 = int64(tVal)
	case float32:
		i64 = int64(tVal)
	case float64:
		i64 = int64(tVal)
	case time.Duration:
		i64 = int64(tVal)
	case Int64able: // eg: json.Number
		i64, err = tVal.Int64()
	default:
		err = ErrConvType
	}
	return
}

// ToFloat convert value to float64, return error on failed
func toFloat(in any) (float64, error) { return toFloatWith(in) }

// ToFloatWith try to convert value to float64. can with some option func, more see ConvOption.
func toFloatWith(in any) (f64 float64, err error) {
	if in == nil {
		return 0, nil
	}

	switch tVal := in.(type) {
	case string:
		f64, err = strconv.ParseFloat(strings.TrimSpace(tVal), 64)
	case int:
		f64 = float64(tVal)
	case int8:
		f64 = float64(tVal)
	case int16:
		f64 = float64(tVal)
	case int32:
		f64 = float64(tVal)
	case int64:
		f64 = float64(tVal)
	case uint:
		f64 = float64(tVal)
	case uint8:
		f64 = float64(tVal)
	case uint16:
		f64 = float64(tVal)
	case uint32:
		f64 = float64(tVal)
	case uint64:
		f64 = float64(tVal)
	case float32:
		f64 = float64(tVal)
	case float64:
		f64 = tVal
	case *float64: // default support float64 ptr type
		f64 = *tVal
	case time.Duration:
		f64 = float64(tVal)
	case Float64able: // eg: json.Number
		f64, err = tVal.Float64()
	default:
		err = ErrConvType
	}
	return
}

//
// region array utils
// - from github.com/gookit/goutil/arrutil
//

// AnyToSlice convert any(allow: array,slice) to []any
func anyToSlice(sl any) (ls []any, err error) {
	rfKeys := reflect.ValueOf(sl)
	if rfKeys.Kind() != reflect.Slice && rfKeys.Kind() != reflect.Array {
		return nil, errors.New("the input param type is invalid")
	}

	for i := 0; i < rfKeys.Len(); i++ {
		ls = append(ls, rfKeys.Index(i).Interface())
	}
	return
}

// IsSubList check given values is sub-list of sample list.
func arrContainsAll[T ScalarType](list, values []T) bool {
	for _, value := range values {
		if !inArray(value, list) {
			return false
		}
	}
	return true
}

// In check the given value whether in the list
func inArray[T ScalarType](value T, list []T) bool {
	for _, elem := range list {
		if elem == value {
			return true
		}
	}
	return false
}

//
// region map utils
// - from github.com/gookit/goutil/maputil
//

// HasKey check of the given map.
func mapHasKey(mp, key any) (ok bool) {
	rftVal := reflect.Indirect(reflect.ValueOf(mp))
	if rftVal.Kind() != reflect.Map {
		return
	}

	for _, keyRv := range rftVal.MapKeys() {
		if reflectIsEqual(keyRv.Interface(), key) {
			return true
		}
	}
	return
}

// HasOneKey check of the given map. return the first exist key
func mapHasOneKey(mp any, keys ...any) (ok bool, key any) {
	rftVal := reflect.Indirect(reflect.ValueOf(mp))
	if rftVal.Kind() != reflect.Map {
		return
	}

	for _, key = range keys {
		for _, keyRv := range rftVal.MapKeys() {
			if reflectIsEqual(keyRv.Interface(), key) {
				return true, key
			}
		}
	}

	return false, nil
}

// HasAllKeys check of the given map. return the first not exist key
func mapHasAllKeys(mp any, keys ...any) (ok bool, noKey any) {
	rftVal := reflect.Indirect(reflect.ValueOf(mp))
	if rftVal.Kind() != reflect.Map {
		return
	}

	for _, key := range keys {
		var exist bool
		for _, keyRv := range rftVal.MapKeys() {
			if reflectIsEqual(keyRv.Interface(), key) {
				exist = true
				break
			}
		}

		if !exist {
			return false, key
		}
	}

	return true, nil
}

//
// region reflect utils
// - from github.com/gookit/goutil/reflects
//

// IsFunc value
func reflectIsFunc(val any) bool {
	if val == nil {
		return false
	}
	return reflect.TypeOf(val).Kind() == reflect.Func
}

// Elem returns the value that the interface v contains
// or that the pointer v points to. otherwise, will return self
func reflectElem(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		return v.Elem()
	}
	return v
}

// IsEqual determines if two objects are considered equal.
//
// TIP: cannot compare a function type
func reflectIsEqual(src, dst any) bool {
	if src == nil || dst == nil {
		return src == dst
	}

	bs1, ok := src.([]byte)
	if !ok {
		return reflect.DeepEqual(src, dst)
	}

	bs2, ok := dst.([]byte)
	if !ok {
		return false
	}

	if bs1 == nil || bs2 == nil {
		return bs1 == nil && bs2 == nil
	}
	return bytes.Equal(bs1, bs2)
}

// IsEmpty reflect value check. if is ptr, check if is nil
func reflectIsEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr, reflect.Func:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}

// Len get reflect value length. allow: intX, uintX, floatX, string, map, array, chan, slice.
//
// Note: (u)intX use width. float to string then calc len.
func reflectLen(v reflect.Value) int {
	v = reflect.Indirect(v)

	// (u)int use width.
	switch v.Kind() {
	case reflect.String:
		return len([]rune(v.String()))
	case reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return v.Len()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return len(strconv.FormatInt(int64(v.Uint()), 10))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return len(strconv.FormatInt(v.Int(), 10))
	case reflect.Float32, reflect.Float64:
		return len(fmt.Sprint(v.Interface()))
	default:
		return -1 // cannot get length
	}
}

//
// region string utils
// - from github.com/gookit/goutil/strutil
//

// check is number: int or float
var numReg = regexp.MustCompile(`^[-+]?\d*\.?\d+$`)

// IsNumeric returns true if the given string is a numeric, otherwise false.
func strIsNumeric(s string) bool { return numReg.MatchString(s) }

// ContainsAll given string should contain all substrings. alias of HasAllSubs()
func strContainsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

//
// region fs utils
// - from github.com/gookit/goutil/fsutil
//

// IsDir reports whether the named directory exists.
func isDir(path string) bool {
	if path == "" || len(path) > 468 {
		return false
	}

	if fi, err := os.Stat(path); err == nil {
		return fi.IsDir()
	}
	return false
}

// IsFile reports whether the named file or directory exists.
func isFile(path string) bool {
	if path == "" || len(path) > 468 {
		return false
	}

	if fi, err := os.Stat(path); err == nil {
		return !fi.IsDir()
	}
	return false
}
