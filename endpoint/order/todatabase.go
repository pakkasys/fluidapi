package order

import (
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
)

var DirectionDatabaseTranslations = map[OrderDirection]util.OrderDirection{
	DIRECTION_ASC:        util.OrderAsc,
	DIRECTION_ASCENDING:  util.OrderAsc,
	DIRECTION_DESC:       util.OrderDesc,
	DIRECTION_DESCENDING: util.OrderDesc,
}

// ValidateAndTranslateToDBOrders validates and translates the provided
// orders into database orders.
// It also returns an error if any of the orders are invalid.
//
//   - orders: The list of orders to validate and translate.
//   - allowedOrderFields: The list of allowed order fields.
//   - apiToDatabaseFieldTranslation: The mapping of API field names to database
//     field names.
func ValidateAndTranslateToDBOrders(
	orders []Order,
	allowedOrderFields []string,
	apiToDatabaseFieldTranslation map[string]dbfield.DBField,
) ([]util.Order, error) {
	newOrders, err := ValidateAndDeduplicateOrders(
		orders,
		allowedOrderFields,
	)
	if err != nil {
		return nil, err
	}

	dbOrders, err := ToDBOrders(newOrders, apiToDatabaseFieldTranslation)
	if err != nil {
		return nil, err
	}

	return dbOrders, nil
}

// ToDBOrders translates the provided orders into database orders.
// It returns an error if any of the orders are invalid.
//
//   - orders: The list of orders to translate.
//   - fieldTranslations: The mapping of API field names to database field
//     names.
func ToDBOrders(
	orders []Order,
	fieldTranslations map[string]dbfield.DBField,
) ([]util.Order, error) {
	newOrders := []util.Order{}

	for i := range orders {
		order := orders[i]

		translatedField := fieldTranslations[order.Field]

		// Translate column
		dbColumn := translatedField.Column
		if dbColumn == "" {
			return nil, InvalidOrderFieldError.WithData(
				InvalidOrderFieldErrorData{
					Field: order.Field,
				},
			)
		}
		order.Field = dbColumn

		newOrders = append(
			newOrders,
			util.Order{
				Table:     translatedField.Table,
				Field:     order.Field,
				Direction: DirectionDatabaseTranslations[order.Direction],
			},
		)
	}

	return newOrders, nil
}
