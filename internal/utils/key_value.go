package utils

import "strings"

type KeyValueEntry[T any] struct {
	Key   string
	Value T
}

func (a KeyValueEntry[T]) Compare(b KeyValueEntry[T]) int {
	return strings.Compare(a.Key, b.Key)
}
