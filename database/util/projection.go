package util

import "fmt"

type Projection struct {
	Table  string
	Column string
	Alias  string
}

func NewProjection(table string, column string, alias string) Projection {
	return Projection{
		Table:  table,
		Column: column,
		Alias:  alias,
	}
}

func (c *Projection) String() string {
	return fmt.Sprintf("`%s`.`%s` AS `%s`", c.Table, c.Column, c.Alias)
}
