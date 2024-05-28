// Package pointer helps getting a pointer to a value, given the literal value.
// Provided benchmarks show that these helpers do not introduce significant latency
// as long as they're used with literal values (as opposed to variables).
package pointer

import (
	"encoding/json"
	"time"

	"golang.org/x/exp/constraints"
)

// To returns a pointer from any types
func To[K any](input K) *K {
	return &input
}

// From returns the value from any pointer type
func From[K any](input *K) K {
	if input == nil {
		var defaultInput K
		return defaultInput
	}
	return *input
}

// Bool returns a pointer to a specified bool.
func Bool(v bool) *bool { return &v }

// Int returns a pointer to an int.
func Int(v int) *int { return &v }

// Int32 returns a pointer to an int32.
func Int32(v int32) *int32 { return &v }

// Int64 returns a pointer to an int64.
func Int64(v int64) *int64 { return &v }

// Int64Value returns the value of the int64 pointer passed in or 0 if the pointer is nil.
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// String returns a pointer to a string.
func String(v string) *string { return &v }

// StringOrNil returns a pointer to a specified string.nil if blank.
func StringOrNil(v string) *string {
	if v == "" {
		return nil
	}
	return String(v)
}

// StringValue returns the value of the string pointer passed in or "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// Uint8 returns a pointer to a uint8.
func Uint8(v uint8) *uint8 { return &v }

// Uint32 returns a pointer to a uint32.
func Uint32(v uint32) *uint32 { return &v }

// Uint64 returns a pointer to a uint64.
func Uint64(v uint64) *uint64 { return &v }

// Uint returns a pointer to a uint.
func Uint(v uint) *uint { return &v }

// Time returns a pointer to a specified time.Time.
func Time(v time.Time) *time.Time { return &v }

// Float64 returns a pointer to a float64.
func Float64(v float64) *float64 { return &v }

// Float32 returns a pointer to a float32.
func Float32(v float32) *float32 { return &v }

// TimeOrNil returns a pointer to a specified time.Time if not zero.
func TimeOrNil(v time.Time) *time.Time {
	if v.IsZero() {
		return nil
	}
	return Time(v)
}

// JSONRawMessage returns a pointer to the provided json raw message.
func JSONRawMessage(v json.RawMessage) *json.RawMessage { return &v }

// Cast allows casting between integer types like *uint32 => *uint
func Cast[A, B constraints.Integer | constraints.Float](a *A) *B {
	if a == nil {
		return nil
	}

	return To(B(*a))
}
