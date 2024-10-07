package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	sourceTag = "source"
	jsonTag   = "json"

	requestURL     = "url"
	requestBody    = "body"
	requestHeaders = "headers"
	requestCookies = "cookies"
)

// URLEncoder interface for encoding URL values.
type URLEncoder interface {
	EncodeURL(data map[string]any) (url.Values, error)
}

func processAndSend[Payload any](
	client *http.Client,
	host string,
	url string,
	method string,
	input *RequestData,
	urlEncoder URLEncoder,
) (*SendResult[Payload], error) {
	if input.Body != nil && method == http.MethodGet {
		return nil, fmt.Errorf("body cannot be set for GET requests")
	}

	bodyReader, err := marshalBody(input.Body)
	if err != nil {
		return nil, err
	}

	constructedURL, err := constructURL(
		host,
		url,
		input.URLParameters,
		urlEncoder,
	)
	if err != nil {
		return nil, err
	}

	req, err := createRequest(
		method,
		*constructedURL,
		bodyReader,
		input.Headers,
		input.Cookies,
	)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	output, err := responseToPayload(resp, new(Payload))
	if err != nil {
		return nil, err
	}

	return &SendResult[Payload]{resp, output}, nil
}

func createRequest(
	method string,
	fullURL string,
	bodyReader io.Reader,
	headers map[string]string,
	cookies []http.Cookie,
) (*http.Request, error) {
	var body io.Reader
	if bodyReader != nil {
		body = bodyReader
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}

	return req, nil
}

func marshalBody(body any) (*bytes.Reader, error) {
	if body == nil {
		return bytes.NewReader(nil), nil
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bodyBytes), nil
}

func constructURL(
	host string,
	path string,
	urlParameters map[string]any,
	urlEncoder URLEncoder,
) (*string, error) {
	url := fmt.Sprintf("%s%s", host, path)

	// Return the base URL if there are no URL parameters
	if len(urlParameters) == 0 {
		return &url, nil
	}

	params, err := toURLParamString(urlParameters, urlEncoder)
	if err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf("%s?%s", url, *params)
	return &fullURL, nil
}

func toURLParamString(
	input map[string]any,
	urlEncoder URLEncoder,
) (*string, error) {
	if input == nil {
		return nil, fmt.Errorf("input map is nil")
	}

	values, err := urlEncoder.EncodeURL(input)
	if err != nil {
		return nil, err
	}

	urlParams := values.Encode()
	return &urlParams, nil
}

func responseToPayload[T any](r *http.Response, output *T) (*T, error) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, output); err != nil {
		return nil, fmt.Errorf(
			"JSON unmarshal error: %v, body: %s", err, string(body),
		)
	}

	return output, nil
}

func parseInput(method string, input any) (*ParsedInput, error) {
	if input == nil {
		return nil, fmt.Errorf("parsed input is nil")
	}

	// Initialize maps and slices for parsed data
	headers := make(map[string]string)
	cookies := make([]http.Cookie, 0)
	urlParameters := make(map[string]any)
	body := make(map[string]any)

	defaultPlacement := determineDefaultPlacement(method)

	// Extract values from the input struct and process them based on their tags
	inputVal := reflect.ValueOf(input).Elem()
	inputType := inputVal.Type()

	for i := 0; i < inputVal.NumField(); i++ {
		field := inputVal.Field(i)
		fieldInfo := inputType.Field(i)
		updatedCookies, err := processField(
			field,
			fieldInfo,
			defaultPlacement,
			headers,
			cookies,
			urlParameters,
			body,
		)
		if err != nil {
			return nil, err
		}
		cookies = updatedCookies
	}

	return &ParsedInput{
		Headers:       headers,
		Cookies:       cookies,
		URLParameters: urlParameters,
		Body:          body,
	}, nil
}

// determineDefaultPlacement determines the default placement of fields based on
// the HTTP method.
func determineDefaultPlacement(method string) string {
	switch method {
	case http.MethodGet:
		return requestURL
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return requestBody
	default:
		return requestBody
	}
}

// processField processes a single struct field and updates the appropriate map
// or slice.
func processField(
	field reflect.Value,
	fieldInfo reflect.StructField,
	defaultPlacement string,
	headers map[string]string,
	cookies []http.Cookie,
	urlParameters map[string]any,
	body map[string]any,
) ([]http.Cookie, error) {
	// Determine field value placement (e.g. URL, body, headers, cookies)
	placement := fieldInfo.Tag.Get(sourceTag)
	if placement == "" {
		placement = defaultPlacement
	}

	// Place the field value in the appropriate map or slice
	return placeFieldValue(
		placement,
		determineFieldName(
			fieldInfo.Tag.Get(jsonTag),
			fieldInfo.Name,
		),
		extractFieldValue(field),
		headers,
		cookies,
		urlParameters,
		body,
	)
}

// determineFieldName determines the field name to use based on the JSON tag.
func determineFieldName(jsonTag string, fieldName string) string {
	jsonFieldName := strings.Split(jsonTag, ",")[0]
	if jsonFieldName == "" {
		jsonFieldName = fieldName
	}
	return jsonFieldName
}

// extractFieldValue extracts the value of a field based on its kind.
func extractFieldValue(field reflect.Value) any {
	switch field.Kind() {
	case reflect.Bool:
		return field.Bool()
	case reflect.String:
		return field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint()
	case reflect.Float32, reflect.Float64:
		return field.Float()
	default:
		return field.Interface()
	}
}

// placeFieldValue places the extracted field value into the appropriate map or
// slice.
func placeFieldValue(
	placement string,
	jsonFieldName string,
	value any,
	headers map[string]string,
	cookies []http.Cookie,
	urlParameters map[string]any,
	body map[string]any,
) ([]http.Cookie, error) {
	switch placement {
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
		return cookies, fmt.Errorf("invalid source tag: %s", placement)
	}
	return cookies, nil
}
