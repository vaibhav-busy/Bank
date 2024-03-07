package models

type Bank struct {
	ID       uint64    `pg:",pk"`
	Name     string    `pg:",unique,notnull"`
	Branches []*Branch `pg:"rel:has-many"`
}
