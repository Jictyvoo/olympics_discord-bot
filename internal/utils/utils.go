package utils

import (
	"os"
	"time"
)

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

type numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func AbsoluteNum[T numeric](value T) T {
	if value < 0 {
		return -value
	}

	return value
}

func EnsureTime(value *time.Time, duration time.Duration) time.Time {
	if value == nil {
		value = new(time.Time)
	}
	if !value.IsZero() {
		return *value
	}

	now := time.Now()
	*value = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).
		Add(duration)

	return *value
}
