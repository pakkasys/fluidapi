package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	clienturl "github.com/PakkaSys/fluidapi/core/client/url"
)

const (
	sourceTag = "source"
	jsonTag   = "json"

	requestURL     = "url"
	requestBody    = "body"
	requestHeaders = "headers"
	requestCookies = "cookies"
)

func ResponseToPayload[T any](response *http.Response, output *T) (*T, error) {
	body, err := ReadBody(response)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, output); err != nil {
		return nil, fmt.Errorf(
			"JSON unmarshal error: %v, body: %s",
			err,
			string(body),
		)
	}

	return output, nil
}

func ReadBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func toURLParamString(input any) *string {
	v := reflect.ValueOf(input)

	if v.Kind() != reflect.Struct &&
		v.Kind() != reflect.Slice &&
		v.Kind() != reflect.Map {

		panic("no valid input provided for URL parameter conversion")
	}

	values := url.Values{}
	kind := v.Kind()

	if kind == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldTag := v.Type().Field(i).Tag.Get("json")

			// Skip if the json tag is not set or empty
			if fieldTag == "" {
				panic(fmt.Sprintf(
					"cannot encode field %q because it has no json tag",
					v.Type().Field(i).Name,
				))
			}

			clienturl.EncodeStructToURL(
				&values,
				fieldTag,
				field,
			)
		}
	} else if kind == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			elemTag := fmt.Sprintf("elem%d", i)
			clienturl.EncodeStructToURL(
				&values,
				elemTag,
				elem,
			)
		}
	} else if kind == reflect.Map {
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			clienturl.EncodeStructToURL(
				&values,
				key.String(),
				value,
			)
		}
	}

	urlParams := values.Encode()

	return &urlParams
}

func MarshalBody(body any) (*bytes.Reader, error) {
	if body == nil {
		return bytes.NewReader(nil), nil
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bodyBytes), nil
}

func ConstructURL(
	host string,
	path string,
	urlParameters map[string]any,
) string {
	url := fmt.Sprintf("%s%s", host, path)

	if len(urlParameters) != 0 {
		return fmt.Sprintf(
			"%s?%s",
			url,
			*toURLParamString(urlParameters),
		)
	}

	return url
}

func ParseInput[Input any](method string, input *Input) (
	map[string]string,
	[]http.Cookie,
	map[string]any,
	map[string]any,
	error,
) {
	if input == nil {
		return nil, nil, nil, nil, fmt.Errorf("input is nil")
	}

	headers := make(map[string]string)
	cookies := make([]http.Cookie, 0)
	urlParameters := make(map[string]any)
	body := make(map[string]any)

	var defaultPlacement string
	switch method {
	case http.MethodGet:
		defaultPlacement = requestURL
	case http.MethodPost:
		defaultPlacement = requestBody
	case http.MethodPut:
		defaultPlacement = requestBody
	case http.MethodPatch:
		defaultPlacement = requestBody
	case http.MethodDelete:
		defaultPlacement = requestBody
	default:
		defaultPlacement = requestBody
	}

	inputVal := reflect.ValueOf(input).Elem()
	inputType := inputVal.Type()

	for i := 0; i < inputVal.NumField(); i++ {
		field := inputVal.Field(i)
		fieldInfo := inputType.Field(i)
		jsonTag := fieldInfo.Tag.Get(jsonTag)

		tag := fieldInfo.Tag.Get(sourceTag)
		if tag == "" {
			tag = defaultPlacement
		}

		jsonFieldName := strings.Split(jsonTag, ",")[0]
		if jsonFieldName == "" {
			jsonFieldName = fieldInfo.Name
		}

		var value any
		switch field.Kind() {
		case reflect.Bool:
			value = field.Bool()
		case reflect.String:
			value = field.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = field.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = field.Uint()
		case reflect.Float32, reflect.Float64:
			value = field.Float()
		default:
			value = field.Interface()
		}

		switch tag {
		case requestURL:
			urlParameters[jsonFieldName] = value
		case requestBody:
			body[jsonFieldName] = value
		case requestHeaders:
			headers[jsonFieldName] = fmt.Sprintf("%v", value)
		case requestCookies:
			cookies = append(
				cookies,
				http.Cookie{
					Name:  jsonFieldName,
					Value: fmt.Sprintf("%v", value),
				},
			)
		default:
			return nil, nil, nil, nil, fmt.Errorf("invalid source tag: %s", tag)
		}
	}

	return headers, cookies, urlParameters, body, nil
}
