package util

import (
	"fmt"
	"strings"
)

type Projection struct {
	Table  string
	Column string
	Alias  string
}

func (c *Projection) String() string {
	builder := strings.Builder{}

	if c.Table == "" {
		builder.WriteString(fmt.Sprintf("`%s`", c.Column))
	} else {
		builder.WriteString(fmt.Sprintf("`%s`.`%s`", c.Table, c.Column))
	}

	if c.Alias != "" {
		builder.WriteString(fmt.Sprintf(" AS `%s`", c.Alias))
	}

	return builder.String()
}
