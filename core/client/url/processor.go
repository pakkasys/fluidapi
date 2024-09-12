package url

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Decodes URL from following syntax:
// someKey=value
// someStruct.field=value
// someSlice[0].field=value
func DecodeURL(values url.Values) map[string]any {
	urlData := make(map[string]any)
	for k, v := range values {
		setNestedMapValue(urlData, k, v[0])
	}
	return urlData
}

// Encodes URL into following syntax:
// someKey=value
// someStruct.field=value
// someSlice[0].field=value
// Every struct field must have a json tag or else it will panic.
func EncodeStructToURL(values *url.Values, fieldTag string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		// Dereference the pointer and recursively call the function
		// But first, check if it's not nil to avoid panic
		if !v.IsNil() {
			EncodeStructToURL(
				values,
				fieldTag,
				v.Elem(),
			)
		}
	case reflect.String:
		values.Set(fieldTag, v.String())
	case reflect.Int, reflect.Int32, reflect.Int64:
		values.Set(fieldTag, fmt.Sprintf("%d", v.Int()))
	case reflect.Bool:
		boolVal := v.Bool()
		values.Set(fieldTag, strconv.FormatBool(boolVal))
	case reflect.Slice:
		for j := 0; j < v.Len(); j++ {
			sliceElem := v.Index(j)
			newFieldTag := fmt.Sprintf("%s[%d]", fieldTag, j)
			EncodeStructToURL(
				values,
				newFieldTag,
				sliceElem,
			)
		}
	case reflect.Struct:
		// Special case for time.Time
		if v.Type() == reflect.TypeOf(time.Time{}) {
			values.Set(
				fieldTag,
				v.Interface().(time.Time).Format(time.RFC3339),
			)
			return
		}
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := v.Type().Field(i)

			// If it's an embedded struct, we recursively call the function
			// without appending a new tag
			if fieldType.Anonymous {
				EncodeStructToURL(values, fieldTag, field)
				continue
			}

			newFieldTag := fieldType.Tag.Get("json")

			// Skip if the json tag is not set or empty
			if newFieldTag == "-" || newFieldTag == "" {
				panic(fmt.Sprintf(
					"cannot encode field %q because it has no json tag",
					fieldType.Name,
				))
			}

			// If the field tag doesn't already have a prefix, do not add a dot
			if fieldTag != "" {
				newFieldTag = fieldTag + "." + newFieldTag
			}

			EncodeStructToURL(values, newFieldTag, field)
		}
	default:
		panic(fmt.Sprintf(
			"value type not supported by URL encoding: %s",
			v.Kind(),
		))
	}
}

// Helper function to set a value in a nested map or slice based on a compound key.
func setNestedMapValue(m map[string]any, key string, value any) {
	// Pattern to match array-like indices e.g., someSlice[0]
	r := regexp.MustCompile(`(\w+)\[(\d+)\]`)
	parts := strings.Split(key, ".")

	var current any = m
	for i, part := range parts {
		if i == len(parts)-1 {
			// Final part where the value must be set
			if sliceIndex := r.FindStringSubmatch(part); sliceIndex != nil {
				// It's a slice index
				sliceName, index := sliceIndex[1], sliceIndex[2]
				idx, err := strconv.Atoi(index)
				if err != nil {
					panic(fmt.Sprintf("invalid index: %s", index))
				}
				currentMap := current.(map[string]any)
				if _, ok := currentMap[sliceName]; !ok {
					// Make sure the slice can accommodate the index
					currentMap[sliceName] = make([]any, idx+1)
				}
				slice := currentMap[sliceName].([]any)
				if idx >= len(slice) {
					// Extend the slice if the index is out of bounds
					newSlice := make([]any, idx+1)
					copy(newSlice, slice)
					slice = newSlice
				}
				slice[idx] = value
				currentMap[sliceName] = slice
			} else {
				// Regular key
				current.(map[string]any)[part] = value
			}
			return
		}

		// Intermediate parts, ensure map exists
		if sliceIndex := r.FindStringSubmatch(part); sliceIndex != nil {
			// It is a slice index
			sliceName, index := sliceIndex[1], sliceIndex[2]
			idx, err := strconv.Atoi(index)
			if err != nil {
				panic(fmt.Sprintf("invalid index: %s", index))
			}
			currentMap := current.(map[string]any)
			if _, ok := currentMap[sliceName]; !ok {
				// Initial slice
				currentMap[sliceName] = make([]any, idx+1)
			}
			slice := currentMap[sliceName].([]any)
			if idx >= len(slice) {
				// Extend the slice if the index is out of bounds
				newSlice := make([]any, idx+1)
				copy(newSlice, slice)
				slice = newSlice
			}
			if _, ok := slice[idx].(map[string]any); !ok {
				slice[idx] = make(map[string]any)
			}
			current = slice[idx]
		} else {
			// Regular map
			currentMap := current.(map[string]any)
			if _, ok := currentMap[part]; !ok {
				currentMap[part] = make(map[string]any)
			}
			current = currentMap[part]
		}
	}
}
