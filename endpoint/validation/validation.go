package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

const (
	errorFmt = "Field %q validation failed on rule %q"
)

// Validation is a validator using the validator/v10 package.
type Validation struct {
	validate *validator.Validate
}

// NewValidation creates a new ValidatorService instance.
func NewValidation() *Validation {
	return &Validation{
		validate: validator.New(),
	}
}

// ValidateStruct validates a struct and returns an error if validation fails.
func (vs *Validation) ValidateStruct(obj any) error {
	return vs.validate.Struct(obj)
}

// ValidateVariable validates a variable against a rule and returns an error if
// validation fails.
func (v *Validation) ValidateVariable(
	fieldName string,
	obj any,
	rule string,
) error {
	err := v.validate.Var(obj, rule)
	if err != nil {
		return v.fieldValidationError(fieldName, err)
	}
	return nil
}

// ValidateAndReturnErrors validates the given object and returns a slice of
// error messages.
func (vs *Validation) GetErrorStrings(err error) []string {
	if err == nil {
		return nil
	}
	validationErrors := parseValidationErrors(err)
	return validationErrors
}

func (v *Validation) fieldValidationError(fieldName string, err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	if len(validationErrors) == 0 {
		return nil
	}

	return fmt.Errorf(errorFmt, fieldName, validationErrors[0].Tag())
}

func parseValidationErrors(err error) []string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Check if it's a custom formatted error with field name
		return []string{err.Error()}
	}

	var errorMessages []string
	for _, err := range validationErrors {
		errorMessages = append(
			errorMessages,
			fmt.Sprintf(errorFmt, err.Field(), err.Tag()),
		)
	}

	return errorMessages
}
