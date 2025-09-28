package dto

import (
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ct.Time = time.Time{}
		return nil
	}

	s := string(b[1 : len(b)-1])
	if s == "" {
		ct.Time = time.Time{}
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	ct.Time = t

	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.Time.Format("2006-01-02") + `"`), nil
}
