package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var (
	_ json.Marshaler   = (*Duration)(nil)
	_ json.Unmarshaler = (*Duration)(nil)
)

type Duration struct{ time.Duration }

func (d *Duration) MarshalJSON() ([]byte, error) {
	ms := int64(d.Duration / time.Millisecond)
	raw := fmt.Sprintf("%d", ms)
	return []byte(raw), nil
}

func (d *Duration) UnmarshalJSON(raw []byte) error {
	ms, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return err
	}

	*d = Duration{time.Duration(ms) * time.Millisecond}
	return nil
}

var (
	_ json.Marshaler   = (*Duration)(nil)
	_ json.Unmarshaler = (*Duration)(nil)
)

type Time struct{ time.Time }

func (t *Time) MarshalJSON() ([]byte, error) {
	ms := t.Time.UnixMilli()
	raw := fmt.Sprintf("%d", ms)
	return []byte(raw), nil
}

func (t *Time) UnmarshalJSON(raw []byte) error {
	ms, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return err
	}

	*t = Time{time.UnixMilli(ms)}
	return nil
}
