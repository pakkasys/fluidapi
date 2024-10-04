package order

import (
	"slices"
)

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

		field := order.GetField()
		if !addedFields[field] {
			newOrders = append(newOrders, order)
			addedFields[field] = true
		}
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
