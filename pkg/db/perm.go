package db

import (
	"database/sql/driver"
	"errors"
)

type Perm uint8

const (
	PermDeny Perm = iota
	PermIgnore
	PermAllow
)

// ToBool returns true if and only if this Perm is equal to PermAllow
func (p Perm) ToBool() bool {
	if p == PermAllow {
		return true
	}
	return false
}

// Value - Implement the database/sql Valuer interface
func (p Perm) Value() (driver.Value, error) {
	return int64(p), nil
}

// Scan - Implement the database/sql Scanner interface
func (p *Perm) Scan(value interface{}) error {
	if value == nil {
		*p = PermIgnore
		return nil
	}
	if bv, err := driver.Int32.ConvertValue(value); err == nil {
		if v, ok := bv.(int64); ok {
			*p = Perm(v)
			return nil
		}
	}
	return errors.New("failed to scan Perm")
}
