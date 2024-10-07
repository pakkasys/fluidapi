package update

import (
	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/database/entity"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
)

type InvalidDatabaseUpdateTranslationErrorData struct {
	Field string `json:"field"`
}

var InvalidDatabaseUpdateTranslationError = api.NewError[InvalidDatabaseUpdateTranslationErrorData]("INVALID_DATABASE_UPDATE_TRANSLATION")

// Update represents a data update with a field and a value.
type Update struct {
	Field string // The field to be updated
	Value any    // The new value for the field
}

// ToDBUpdates translates a list of updates to a database update list
// and returns an error if the translation fails.
//
// Parameters:
// - updates: The list of updates to translate.
// - apiToDBFieldMap: The mapping of API field names to database field names.
//
// Returns:
// - A list of database entity updates.
// - An error if any field translation fails.
func ToDBUpdates(
	updates []Update,
	apiToDBFieldMap map[string]dbfield.DBField,
) ([]entity.Update, error) {
	var dbUpdates []entity.Update

	for i := range updates {
		matchedUpdate := updates[i]

		// Translate the field
		dbField, ok := apiToDBFieldMap[matchedUpdate.Field]
		if !ok {
			return nil, InvalidDatabaseUpdateTranslationError.WithData(
				InvalidDatabaseUpdateTranslationErrorData{
					Field: matchedUpdate.Field,
				},
			)
		}

		dbUpdates = append(
			dbUpdates,
			entity.Update{
				Field: dbField.Column,
				Value: matchedUpdate.Value,
			},
		)
	}

	return dbUpdates, nil
}
