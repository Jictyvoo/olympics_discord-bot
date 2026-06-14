package repoinj

import (
	"fmt"

	"github.com/wrapped-owls/goremy-di/remy"
)

type InjectionKind uint8

const (
	InjectionFactory InjectionKind = iota
	InjectionSingleton
	InjectionLazySingleton
)

func ensureImplements[T, I any]() error {
	var t T

	_, ok := any(t).(I)
	if !ok {
		return fmt.Errorf("%T does not implement requested interface", t)
	}

	return nil
}

func RegisterAlias[T, I any](container remy.Injector, kind InjectionKind) error {
	bindFuncInterface := remy.Factory[I]
	switch kind {
	case InjectionFactory:
		// already set
	case InjectionSingleton:
		bindFuncInterface = remy.Singleton[I]
	case InjectionLazySingleton:
		bindFuncInterface = remy.LazySingleton[I]
	default:
		return fmt.Errorf("invalid injection kind: %v", kind)
	}

	if err := ensureImplements[T, I](); err != nil {
		return err
	}
	remy.RegisterConstructorArgs1(
		container, bindFuncInterface,
		func(val T) I { return any(val).(I) },
	)

	return nil
}

func RegisterAliased[T, A, I any](
	container remy.Injector, kind InjectionKind,
	constructor func(A) T, _ ...*I,
) error {
	bindFuncRaw := remy.Factory[T]
	switch kind {
	case InjectionFactory:
		// already set
	case InjectionSingleton:
		bindFuncRaw = remy.Singleton[T]
	case InjectionLazySingleton:
		bindFuncRaw = remy.LazySingleton[T]
	default:
		return fmt.Errorf("invalid injection kind: %v", kind)
	}

	remy.RegisterConstructorArgs1(
		container, bindFuncRaw, constructor,
	)

	err := RegisterAlias[T, I](container, kind)
	return err
}
