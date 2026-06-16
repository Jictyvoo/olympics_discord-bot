package confloader

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// Binder applies a single environment variable override to a config field.
type Binder interface {
	Bind() error
}

type binder[T any] struct {
	field  *T
	env    string
	parser func(string) (T, error)
}

func (b binder[T]) Bind() error {
	val := os.Getenv(b.env)
	if val == "" {
		return nil
	}
	parsed, err := b.parser(val)
	if err != nil {
		return err
	}
	*b.field = parsed
	return nil
}

//nolint:ireturn // factory returning consumer interface by design
func BindField[T any](field *T, env string, parser func(string) (T, error)) Binder {
	return binder[T]{field: field, env: env, parser: parser}
}

func ParseString(s string) (string, error) { return s, nil }

func ParseBool(s string) (bool, error) { return strconv.ParseBool(s) }

func ParseInt[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}](bitSize int) func(string) (T, error) {
	return func(s string) (T, error) {
		v, err := strconv.ParseInt(s, 10, bitSize)
		return T(v), err
	}
}

func ParseUint[T interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}](bitSize int) func(string) (T, error) {
	return func(s string) (T, error) {
		v, err := strconv.ParseUint(s, 10, bitSize)
		return T(v), err
	}
}

func ParseFloat64(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

func ParseDuration(s string) (time.Duration, error) { return time.ParseDuration(s) }

// BindEnv runs every binder, accumulating all errors rather than stopping at the first.
func BindEnv(binders ...Binder) error {
	errs := make([]error, 0, len(binders))
	for _, b := range binders {
		errs = append(errs, b.Bind())
	}
	return errors.Join(errs...)
}
