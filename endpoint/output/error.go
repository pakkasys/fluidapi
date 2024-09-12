package output

import "github.com/PakkaSys/fluidapi/core/api"

var ERROR_ID = "ERROR"

func Error() *api.Error {
	return &api.Error{
		ID: ERROR_ID,
	}
}
