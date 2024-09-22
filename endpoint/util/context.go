package util

import (
	"context"
	"sync"
)

type DataKey int

var (
	base DataKey = 0
	lock sync.Mutex

	mainDataKey = NewDataKey()
)

type contextData struct {
	data sync.Map
}

// NewDataKey safely increments and returns the next value of base.
func NewDataKey() DataKey {
	lock.Lock()
	defer lock.Unlock()
	base++
	return base
}

// NewContext initializes a new context with an empty contextData map.
func NewContext(fromContext context.Context) context.Context {
	return context.WithValue(fromContext, mainDataKey, &contextData{})
}

// IsContextSet checks if the custom context is set in the context.
func IsContextSet(context context.Context) bool {
	return HasContextValue(context, mainDataKey)
}

// HasContextValue checks if a value exists for the provided key within the
// custom data of the context.
func HasContextValue(context context.Context, payloadKey any) bool {
	cd, ok := getContextData(context)
	if !ok {
		return false
	}

	_, exists := cd.data.Load(payloadKey)
	return exists
}

// GetContextValue tries to retrieve a value from the custom data of the context
// for a given key.
// If the key exists and the value matches the expected type, it returns the
// value. Otherwise, it returns the provided default value.
func GetContextValue[T any](
	context context.Context,
	payloadKey any,
	returnOnNull T,
) T {
	cd, ok := getContextData(context)
	if !ok {
		return returnOnNull
	}

	value, exists := cd.data.Load(payloadKey)
	if !exists {
		return returnOnNull
	}

	typedValue, isType := value.(T)
	if !isType {
		return returnOnNull
	}

	return typedValue
}

// MustGetContextValue fetches a value directly from the custom data of the
// context for a given key.
// This function  panic if the key does not exist or if ther is a type mismatch.
func MustGetContextValue[T any](context context.Context, payloadKey any) T {
	cd, ok := getContextData(context)
	if !ok {
		panic("get context value: no custom context set in request")
	}
	value, exists := cd.data.Load(payloadKey)
	if !exists {
		panic("get context value: key does not exist")
	}

	typedValue, isType := value.(T)
	if !isType {
		panic("get context value: type mismatch")
	}

	return typedValue
}

// SetContextValue sets a value in the custom data of the context for the
// provided key.
func SetContextValue(context context.Context, payloadKey any, payload any) {
	cd, ok := getContextData(context)
	if !ok {
		panic("set context value: no custom context set in request")
	}
	cd.data.Store(payloadKey, payload)
}

// CanSetContextValue checks if the custom data of the context is set.
func CanSetContextValue(context context.Context) bool {
	_, ok := getContextData(context)
	return ok
}

// ClearContextValue clears a value in the custom data of the context for the
// provided key.
func ClearContextValue(context context.Context, payloadKey any) {
	cd, ok := getContextData(context)
	if ok {
		cd.data.Delete(payloadKey)
	}
}

func getContextData(context context.Context) (*contextData, bool) {
	cd, ok := context.Value(mainDataKey).(*contextData)
	return cd, ok && cd != nil
}
