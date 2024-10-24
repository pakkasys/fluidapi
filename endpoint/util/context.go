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
// It is used to create a unique key for storing custom context data.
//
// Returns:
// - The next data key value.
func NewDataKey() DataKey {
	lock.Lock()
	defer lock.Unlock()
	base++
	return base
}

// NewContext initializes a new context with an empty contextData map.
//
// Parameters:
// - fromCtx: The context from which the new context is derived.
//
// Returns:
// - A new context with an initialized custom data map.
func NewContext(fromCtx context.Context) context.Context {
	return context.WithValue(fromCtx, mainDataKey, &contextData{})
}

// IsContextSet checks if the custom context is set in the provided context.
//
// Parameters:
// - ctx: The context to check.
//
// Returns:
// - A boolean value indicating if the custom context is set.
func IsContextSet(ctx context.Context) bool {
	_, ok := getContextData(ctx)
	return ok
}

// HasContextValue checks if a value exists for the provided key within the
// custom data of the context.
//
// Parameters:
// - ctx: The context to check.
// - key: The key for which to check the existence of a value.
//
// Returns:
// - A boolean value indicating if the value exists in the context.
func HasContextValue(ctx context.Context, key any) bool {
	cd, ok := getContextData(ctx)
	if !ok {
		return false
	}

	_, exists := cd.data.Load(key)
	return exists
}

// GetContextValue tries to retrieve a value from the custom data of the context
// for a given key.
// If the key exists and the value matches the expected type, it returns the
// value. Otherwise, it returns the provided default value.
//
// Parameters:
//   - ctx: The context from which to retrieve the value.
//   - key: The key for which to retrieve the value.
//   - returnOnNull: The default value to return if the key does not exist or
//     the type does not match.
//
// Returns:
//   - The value from the context if it exists and matches the expected type,
//     otherwise the default value.
func GetContextValue[T any](
	ctx context.Context,
	key any,
	returnOnNull T,
) T {
	cd, ok := getContextData(ctx)
	if !ok {
		return returnOnNull
	}

	value, exists := cd.data.Load(key)
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
// This function will panic if the key does not exist or if there is a type
// mismatch.
//
// Parameters:
// - ctx: The context from which to fetch the value.
// - key: The key for which to fetch the value.
//
// Returns:
// - The value from the context if it exists and matches the expected type.
//
// Panics:
//   - If the custom context is not set, the key does not exist, or there is a
//     type mismatch.
func MustGetContextValue[T any](ctx context.Context, key any) T {
	cd, ok := getContextData(ctx)
	if !ok {
		panic("get context value: no custom context set in request")
	}
	value, exists := cd.data.Load(key)
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
//
// Parameters:
// - ctx: The context in which to set the value.
// - key: The key for which to set the value.
// - data: The value to set in the context.
//
// Returns:
// - The updated context.
//
// Panics:
// - If the key is nil or if the custom context is not set.
func SetContextValue(ctx context.Context, key any, data any) context.Context {
	if key == nil {
		panic("set context value: key cannot be nil")
	}

	cd, ok := getContextData(ctx)
	if !ok {
		panic("set context value: no custom context set in request")
	}

	cd.data.Store(key, data)

	return ctx
}

// ClearContextValue clears a value in the custom data of the context for the
// provided key.
//
// Parameters:
// - ctx: The context from which to clear the value.
// - key: The key for which to clear the value.
//
// Returns:
// - The updated context.
func ClearContextValue(ctx context.Context, key any) context.Context {
	cd, ok := getContextData(ctx)
	if ok {
		cd.data.Delete(key)
	}

	return ctx
}

// DumpData retrieves the custom context data from the provided context.
//
// Parameters:
// - ctx: The context from which to retrieve the custom data.
//
// Returns:
// - A pointer to the contextData if it exists and is valid, otherwise nil.
func DumpData(ctx context.Context) *contextData {
	cd, _ := getContextData(ctx)
	return cd
}

// getContextData retrieves the custom context data from the provided context.
//
// Parameters:
// - ctx: The context from which to retrieve the custom data.
//
// Returns:
//   - A pointer to the contextData and a boolean indicating if the data exists
//     and is valid.
func getContextData(ctx context.Context) (*contextData, bool) {
	cd, ok := ctx.Value(mainDataKey).(*contextData)
	return cd, ok && cd != nil
}
