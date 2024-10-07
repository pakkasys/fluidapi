package util

import "fmt"

// JoinType represents the type of join
type JoinType string

const (
	JoinTypeInner JoinType = "INNER"
	JoinTypeLeft  JoinType = "LEFT"
	JoinTypeRight JoinType = "RIGHT"
	JoinTypeFull  JoinType = "FULL"
)

// Join represents a database join clause
type Join struct {
	Type    JoinType
	Table   string
	OnLeft  ColumSelector
	OnRight ColumSelector
}

// ColumSelector represents a column selector
type ColumSelector struct {
	Table  string
	Column string
}

// String returns the string representation of the ColumnSelector
func (c *ColumSelector) String() string {
	return fmt.Sprintf("`%s`.`%s`", c.Table, c.Column)
}
