package util

import "fmt"

type JoinType string

const (
	JoinTypeInner JoinType = "INNER"
	JoinTypeLeft  JoinType = "LEFT"
	JoinTypeRight JoinType = "RIGHT"
	JoinTypeFull  JoinType = "FULL"
)

type Join struct {
	Type    JoinType
	Table   string
	OnLeft  ColumSelector
	OnRight ColumSelector
}

func NewJoin(
	typ JoinType,
	table string,
	onLeft ColumSelector,
	onRight ColumSelector,
) Join {
	return Join{
		Type:    typ,
		Table:   table,
		OnLeft:  onLeft,
		OnRight: onRight,
	}
}

type ColumSelector struct {
	Table  string
	Column string
}

func NewColumnSelector(
	table string,
	column string,
) ColumSelector {
	return ColumSelector{
		Table:  table,
		Column: column,
	}
}

func (c *ColumSelector) String() string {
	return fmt.Sprintf("`%s`.`%s`", c.Table, c.Column)
}
