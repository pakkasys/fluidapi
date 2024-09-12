package order

import (
	"fmt"
	"slices"
)

type OrderDirection string

const (
	DIRECTION_ASCENDING       OrderDirection = "ASC"
	DIRECTION_ASCENDING_LONG  OrderDirection = "ASCENDING"
	DIRECTION_DESCENDING      OrderDirection = "DESC"
	DIRECTION_DESCENDING_LONG OrderDirection = "DESCENDING"
)

var Directions []OrderDirection = []OrderDirection{
	DIRECTION_ASCENDING,
	DIRECTION_ASCENDING_LONG,
	DIRECTION_DESCENDING,
	DIRECTION_DESCENDING_LONG,
}

type Order struct {
	Field     string         `json:"field"`
	Direction OrderDirection `json:"direction"`
}

func (s Order) GetField() string {
	return s.Field
}

func Get(field string, direction OrderDirection) *Order {
	if !slices.Contains(Directions, direction) {
		panic(fmt.Sprintf("invalid order direction: %s", direction))
	}

	return &Order{
		Field:     field,
		Direction: direction,
	}
}
