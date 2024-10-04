package order

import (
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
)

var DirectionDatabaseTranslations = map[OrderDirection]util.OrderDirection{
	DIRECTION_ASCENDING:       util.OrderAsc,
	DIRECTION_ASCENDING_LONG:  util.OrderAsc,
	DIRECTION_DESCENDING:      util.OrderDesc,
	DIRECTION_DESCENDING_LONG: util.OrderDesc,
}

func ValidateAndTranslateToDatabaseOrders(
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

	Orders, err := ToDatabaseOrders(
		newOrders,
		apiToDatabaseFieldTranslation,
	)
	if err != nil {
		return nil, err
	}

	return Orders, nil
}

func ToDatabaseOrders(
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
