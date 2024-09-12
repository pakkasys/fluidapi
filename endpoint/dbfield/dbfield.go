package dbfield

type DBField struct {
	Table  string
	Column string
}

func NewDBField(table string, column string) DBField {
	return DBField{
		Table:  table,
		Column: column,
	}
}
