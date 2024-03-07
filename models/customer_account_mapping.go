package models

import "github.com/go-pg/pg/v10/orm"

func init() {
	orm.RegisterTable((*AccountToCustomer)(nil))
}

type AccountToCustomer struct {
	ID         uint64    `pg:",pk"`
	AccountID  uint64    `pg:"on_delete:CASCADE, notnull"`
	CustomerID uint64    `pg:"on_delete:SET NULL, notnull"`
	Account    *Account  `pg:"rel:has-one"`
	Customer   *Customer `pg:"rel:has-one"`
}
