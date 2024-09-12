package runner

import "net/http"

func UpdateInvoker[Input IUpdateInput, Output any](
	specification IUpdateSpecification[Input],
	apiFields APIFields,
	serviceFunc UpdateServiceFunc,
	outputFn func(count int64) *Output,
) func(w http.ResponseWriter, r *http.Request, input *Input) (*Output, error) {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		input *Input,
	) (*Output, error) {
		return UpdateInvoke(
			w,
			r,
			*input,
			specification,
			apiFields,
			serviceFunc,
			outputFn,
		)
	}
}

func DeleteInvoker[Input IDeleteInput, Output any, Entity any](
	specification IDeleteSpecification[Input],
	apiFields APIFields,
	serviceFunc DeleteServiceFunc[Entity],
	outputFn func(count int64) *Output,
) func(w http.ResponseWriter, r *http.Request, input *Input) (*Output, error) {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		input *Input,
	) (*Output, error) {
		return DeleteInvoke(
			w,
			r,
			*input,
			specification,
			apiFields,
			serviceFunc,
			outputFn,
		)
	}
}
