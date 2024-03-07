package models

import "time"

type Customer struct {
	ID       uint64    `pg:",pk"`
	BranchID uint      `pg:"on_delete:CASCADE"`
	Branch   *Branch   `pg:"rel:has-one"`
	Name     string    `pg:",notnull"`
	PAN      string    `pg:",notnull, unique"`
	DOB      time.Time  `pg:",type:date"`
	Age      int
	PhoneNo  string
	Address  string
	Accounts []*Account `pg:"many2many:account_to_customers"`
}
