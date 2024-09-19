package util

import "fmt"

type Projection struct {
	Table  string
	Column string
	Alias  string
}

func (c *Projection) String() string {
	return fmt.Sprintf("`%s`.`%s` AS `%s`", c.Table, c.Column, c.Alias)
}
