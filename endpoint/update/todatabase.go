package update

import (
	"github.com/pakkasys/fluidapi/database/entity"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
)

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
