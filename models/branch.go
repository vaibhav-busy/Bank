package models
import (
	"github.com/google/uuid"
)
type Branch struct {
	ID       uint64 `pg:",pk"`
	Address  string
	BankID   uint64     `pg:"on_delete:CASCADE, notnull, on_update:CASCADE"`
	Bank     *Bank      `pg:"rel:has-one"`
	IfscCode uuid.UUID `pg:"type:uuid"`
	Accounts []*Account `pg:"rel:has-many"`
	Customer []*Customer `pg:"rel:has-many"`
}
