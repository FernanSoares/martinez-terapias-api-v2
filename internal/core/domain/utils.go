package domain

import (
	"encoding/json"
	"time"
)

// TimeOrZero é um wrapper em torno de time.Time para lidar com valores zero
// e formatação de JSON
type TimeOrZero struct {
	time.Time
}

// MarshalJSON implementa a interface json.Marshaler para TimeOrZero
func (t TimeOrZero) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Format("2006-01-02T15:04:05Z07:00"))
}

// UnmarshalJSON implementa a interface json.Unmarshaler para TimeOrZero
func (t *TimeOrZero) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}
	
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		parsedTime, err = time.Parse("2006-01-02", s)
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02T15:04:05", s)
			if err != nil {
				return err
			}
		}
	}
	
	t.Time = parsedTime
	return nil
}
