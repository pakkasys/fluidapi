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
func NewContext(fromCtx context.Context) context.Context {
	return context.WithValue(fromCtx, mainDataKey, &contextData{})
}

// IsContextSet checks if the custom context is set in the context.
func IsContextSet(ctx context.Context) bool {
	return HasContextValue(ctx, mainDataKey)
}

// HasContextValue checks if a value exists for the provided key within the
// custom data of the context.
func HasContextValue(ctx context.Context, payloadKey any) bool {
	cd, ok := getContextData(ctx)
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
	ctx context.Context,
	payloadKey any,
	returnOnNull T,
) T {
	cd, ok := getContextData(ctx)
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
func MustGetContextValue[T any](ctx context.Context, payloadKey any) T {
	cd, ok := getContextData(ctx)
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
func SetContextValue(ctx context.Context, payloadKey any, payload any) {
	cd, ok := getContextData(ctx)
	if !ok {
		panic("set context value: no custom context set in request")
	}
	cd.data.Store(payloadKey, payload)
}

// CanSetContextValue checks if the custom data of the context is set.
func CanSetContextValue(ctx context.Context) bool {
	_, ok := getContextData(ctx)
	return ok
}

// ClearContextValue clears a value in the custom data of the context for the
// provided key.
func ClearContextValue(ctx context.Context, payloadKey any) {
	cd, ok := getContextData(ctx)
	if ok {
		cd.data.Delete(payloadKey)
	}
}

func getContextData(ctx context.Context) (*contextData, bool) {
	cd, ok := ctx.Value(mainDataKey).(*contextData)
	return cd, ok && cd != nil
}
