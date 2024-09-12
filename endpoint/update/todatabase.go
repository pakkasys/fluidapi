package update

import (
	"github.com/PakkaSys/fluidapi/database/entity"
	"github.com/PakkaSys/fluidapi/endpoint/dbfield"
	"github.com/PakkaSys/fluidapi/endpoint/middleware/inputlogic"
)

type MatchedUpdate struct {
	APIUpdate
	InputUpdate
}

type IValidator interface {
	ValidateVariable(fieldName string, obj any, rule string) error
	GetErrorStrings(err error) []string
}

func GetDatabaseUpdatesFromUpdates(
	inputUpdates []InputUpdate,
	allowedUpdates map[string]APIUpdate,
	apiFields map[string]dbfield.DBField,
	validator IValidator,
) ([]entity.Update, error) {
	matchedUpdates, err := MatchAndValidateInputUpdates(
		inputUpdates,
		allowedUpdates,
		validator,
	)
	if err != nil {
		return nil, err
	}

	databaseUpdates, err := ToDatabaseUpdates(apiFields, matchedUpdates)
	if err != nil {
		return nil, err
	}

	return databaseUpdates, nil
}

func MatchAndValidateInputUpdates(
	inputUpdates []InputUpdate,
	allowedUpdates map[string]APIUpdate,
	validator IValidator,
) ([]MatchedUpdate, error) {
	var matchedUpdates []MatchedUpdate

	for i := range inputUpdates {
		inputUpdate := inputUpdates[i]

		// Match the input update to an allowed update
		apiUpdate, ok := allowedUpdates[inputUpdate.Field]
		if !ok {
			return nil, InvalidUpdateFieldError(inputUpdate.Field)
		}

		// Validate the input value
		if err := validator.ValidateVariable(
			inputUpdate.Field,
			inputUpdate.Value,
			apiUpdate.Validation,
		); err != nil {
			return nil, inputlogic.ValidationError(
				validator.GetErrorStrings(err),
			)
		}

		matchedUpdates = append(
			matchedUpdates,
			MatchedUpdate{
				APIUpdate:   apiUpdate,
				InputUpdate: inputUpdate,
			},
		)
	}

	return matchedUpdates, nil
}

func ToDatabaseUpdates(
	apiToDatabaseFieldMap map[string]dbfield.DBField,
	matchedUpdates []MatchedUpdate,
) ([]entity.Update, error) {
	var databaseUpdates []entity.Update

	for i := range matchedUpdates {
		matchedUpdate := matchedUpdates[i]

		// Translate the field
		translatedField, ok := apiToDatabaseFieldMap[matchedUpdate.Field]
		if !ok {
			return nil, InvalidDatabaseUpdateTranslationError(matchedUpdate.Field)
		}

		databaseUpdates = append(
			databaseUpdates,
			*entity.NewUpdate(
				translatedField.Column,
				matchedUpdate.Value,
			),
		)
	}

	return databaseUpdates, nil
}
