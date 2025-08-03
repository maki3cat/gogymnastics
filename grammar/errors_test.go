package grammar

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	// new error
	err1 := errors.New("new error")
	fmt.Println(err1)
	// format error
	err2 := fmt.Errorf("format error: %w", err1)
	fmt.Println(err2)
}

func TestWrapWithFormatAndUnwrap(t *testing.T) {
	// what is the difference between %w and %s?
	// maki: %w wraps and links the error
	err1 := errors.New("error1-inner")
	err2a := fmt.Errorf("error2-outer: %w", err1)
	err2b := fmt.Errorf("error2-outer: %s", err1)
	fmt.Println("error2a", err2a)
	fmt.Println("error2b", err2b)

	// check wrap is stil original type
	fmt.Println("is err2a err1?", errors.Is(err2a, err1))
	fmt.Println("is err2b err1?", errors.Is(err2b, err1))
	fmt.Println("is err2a err2b?", errors.Is(err2a, err2b))

	// check unwrap
	fmt.Println("unwrap err2a", errors.Unwrap(err2a), "Is err1?", errors.Is(errors.Unwrap(err2a), err1))
	fmt.Println("unwrap err2b", errors.Unwrap(err2b), "Is err1?", errors.Is(errors.Unwrap(err2b), err1))
}

func TestWrapWithJoin(t *testing.T) {
	err1 := errors.New("error1-message")
	err2 := errors.New("error2-message")
	err3 := errors.Join(err1, err2)
	fmt.Println("error3:", err3)

	// check the join result
	fmt.Println("is err3 err1?", errors.Is(err3, err1))
	fmt.Println("is err3 err2?", errors.Is(err3, err2))
	fmt.Println("== err3 err1?", err3 == err1)
	fmt.Println("== err3 err2?", err3 == err2)

	// check unwrap
	// ⚠️ Because errors.Unwrap(err3) returns nil for Join.
	unwrapErr3 := errors.Unwrap(err3)
	fmt.Println("unwrap err3", unwrapErr3, "Is err1?", errors.Is(unwrapErr3, err1))
	fmt.Println("unwrap err3", unwrapErr3, "Is err2?", errors.Is(unwrapErr3, err2))
	fmt.Println("unwrap err3", unwrapErr3, "== err1?", unwrapErr3 == err1)
	fmt.Println("unwrap err3", unwrapErr3, "== err2?", unwrapErr3 == err2)
}
