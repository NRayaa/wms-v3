package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Price represents a decimal price that always formats with 2 decimal places
type Price float64

// MarshalJSON ensures the price always has 2 decimal places in JSON
func (p Price) MarshalJSON() ([]byte, error) {
	if p == 0 {
		return []byte("0.00"), nil
	}
	return []byte(fmt.Sprintf("%.2f", p)), nil
}

// UnmarshalJSON parses JSON value into Price
func (p *Price) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	*p = Price(f)
	return nil
}

// Value implements the driver.Valuer interface for database storage
func (p Price) Value() (driver.Value, error) {
	return float64(p), nil
}

// Scan implements the sql.Scanner interface for reading from database
func (p *Price) Scan(value interface{}) error {
	if value == nil {
		*p = 0
		return nil
	}

	// Handle different types from database
	switch v := value.(type) {
	case float64:
		*p = Price(v)
	case float32:
		*p = Price(v)
	case int64:
		*p = Price(v)
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		*p = Price(f)
	case []byte:
		f, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return err
		}
		*p = Price(f)
	case sql.NullFloat64:
		if v.Valid {
			*p = Price(v.Float64)
		} else {
			*p = 0
		}
	default:
		return fmt.Errorf("cannot scan %T into Price", value)
	}
	return nil
}

// String returns the string representation with 2 decimal places
func (p Price) String() string {
	return fmt.Sprintf("%.2f", p)
}
