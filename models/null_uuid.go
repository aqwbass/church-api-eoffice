package models

import (
	"database/sql/driver"

	"github.com/gofrs/uuid"
)

type NullUUID struct {
	UUID  ZeroUUID
	Valid bool
}

func (zu *NullUUID) Scan(src interface{}) error {
	if src == nil {
		zu.UUID, zu.Valid = ZeroUUID(uuid.UUID{}), false
		return nil
	}

	// Delegate to UUID Scan function
	zu.Valid = true
	return zu.UUID.Scan(src)
}

func (zu NullUUID) Value() (driver.Value, error) {
	if !zu.Valid {
		return nil, nil
	}
	// Delegate to UUID Value function
	return zu.UUID.Value()
}
