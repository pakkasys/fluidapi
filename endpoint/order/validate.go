package order

import (
	"slices"
)

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
