package update

import (
	"github.com/pakkasys/fluidapi/database/entity"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
)

// ToDBUpdates translates a list of updates to a database update list
// and returns an error if the translation fails.
//
//   - updates: The list of updates to translate.
//   - apiToDBFieldMap: The mapping of API field names to database field names.
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
