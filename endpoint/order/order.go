package order

import (
	"slices"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/dbfield"
)

type InvalidOrderFieldErrorData struct {
	Field string `json:"field"`
}

var InvalidOrderFieldError = api.NewError[InvalidOrderFieldErrorData]("INVALID_ORDER_FIELD")

type OrderDirection string

const (
	DIRECTION_ASC        OrderDirection = "ASC"
	DIRECTION_ASCENDING  OrderDirection = "ASCENDING"
	DIRECTION_DESC       OrderDirection = "DESC"
	DIRECTION_DESCENDING OrderDirection = "DESCENDING"
)

// Directions is a list of all possible order directions.
var Directions []OrderDirection = []OrderDirection{
	DIRECTION_ASC,
	DIRECTION_ASCENDING,
	DIRECTION_DESC,
	DIRECTION_DESCENDING,
}

var DirectionDatabaseTranslations = map[OrderDirection]util.OrderDirection{
	DIRECTION_ASC:        util.OrderAsc,
	DIRECTION_ASCENDING:  util.OrderAsc,
	DIRECTION_DESC:       util.OrderDesc,
	DIRECTION_DESCENDING: util.OrderDesc,
}

// Order is used to specify the order of the result set.
type Order struct {
	Field     string         `json:"field"`
	Direction OrderDirection `json:"direction"`
}

// ValidateAndDeduplicateOrders validates and deduplicates the provided orders.
// It returns a new list of orders with no duplicates.
// It also returns an error if any of the orders are invalid.
//
//   - orders: The list of orders to validate and deduplicate.
//   - allowedOrderFields: The list of allowed order fields.
func ValidateAndDeduplicateOrders(
	orders []Order,
	allowedOrderFields []string,
) ([]Order, error) {
	newOrders := []Order{}
	addedFields := make(map[string]bool)

	for i := range orders {
		order := orders[i]

		if err := validate(order, allowedOrderFields); err != nil {
			return nil, err
		}

		if !addedFields[order.Field] {
			newOrders = append(newOrders, order)
			addedFields[order.Field] = true
		}
	}

	return newOrders, nil
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

func validate(order Order, allowedOrderFields []string) error {
	// Check that the order direction is valid
	if !slices.Contains(Directions, order.Direction) {
		return InvalidOrderFieldError.WithData(
			InvalidOrderFieldErrorData{
				Field: order.Field,
			},
		)
	}

	// Check that the order field is allowed
	if !slices.Contains(allowedOrderFields, order.Field) {
		return InvalidOrderFieldError.WithData(
			InvalidOrderFieldErrorData{
				Field: order.Field,
			},
		)
	}

	return nil
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
