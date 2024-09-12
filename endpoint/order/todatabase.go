package order

import (
	"github.com/PakkaSys/fluidapi/database/util"
	"github.com/PakkaSys/fluidapi/endpoint/dbfield"
)

var DirectionDatabaseTranslations = map[OrderDirection]util.OrderDirection{
	DIRECTION_ASCENDING:       util.DIRECTION_ASCENDING,
	DIRECTION_ASCENDING_LONG:  util.DIRECTION_ASCENDING,
	DIRECTION_DESCENDING:      util.DIRECTION_DESCENDING,
	DIRECTION_DESCENDING_LONG: util.DIRECTION_DESCENDING,
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
			return nil, InvalidOrderFieldError(order.Field)
		}
		order.Field = dbColumn

		newOrders = append(
			newOrders,
			*util.NewOrder(
				translatedField.Table,
				order.Field,
				DirectionDatabaseTranslations[order.Direction],
			),
		)
	}

	return newOrders, nil
}
