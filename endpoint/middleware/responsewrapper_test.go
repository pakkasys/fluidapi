package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/pakkasys/fluidapi/endpoint/util"
	"github.com/stretchr/testify/assert"
)

func TestSetRequestWrapper(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(util.NewContext(req.Context()))
	mockRequestWrapper, err := util.NewRequestWrapper(req)
	assert.NoError(t, err, "Expected no error when creating request wrapper")

	setRequestWrapper(req, mockRequestWrapper)

	retrievedRequestWrapper := GetRequestWrapper(req)

	assert.NotNil(
		t,
		retrievedRequestWrapper,
		"Request wrapper should not be nil",
	)
	assert.Equal(
		t,
		mockRequestWrapper,
		retrievedRequestWrapper,
		"Request wrapper should match the set request wrapper",
	)
	assert.Equal(
		t,
		req,
		retrievedRequestWrapper.Request,
		"Request inside the wrapper should match the original request",
	)
}
