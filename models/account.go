package models

import "time"

type Account struct {
	ID           uint64 `pg:",pk"`
	BranchID     uint64 `pg:",on_delete:CASCADE, notnull"`
	AccountNo    uint64 `pg:",notnull, unique"`
	Balance      float64
	AccountType  string
	OpeningDate  time.Time      `pg:",type:date"`
	Branch       *Branch        `pg:"rel:has-one"`
	Transactions []*Transaction `pg:"rel:has-many"`
	Customers    []*Customer    `pg:"many2many:account_to_customers"`
}
