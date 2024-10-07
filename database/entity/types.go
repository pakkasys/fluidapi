package entity

type SQLUtil interface {
	CheckDBError(err error) error
}
