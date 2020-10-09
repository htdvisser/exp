package types

import (
	"time"
)

type Time time.Time

const TimeFormat = "2006-01-02T15:04:05.000000"

func (t Time) MarshalText() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat))
	return time.Time(t).AppendFormat(b, TimeFormat), nil
}

func (t Time) String() string {
	b, _ := t.MarshalText()
	return string(b)
}

func (t *Time) UnmarshalText(data []byte) error {
	parsed, err := time.Parse(TimeFormat, string(data))
	if err != nil {
		return err
	}
	*t = Time(parsed)
	return nil
}

func TimePtr(t *time.Time) *Time {
	if t == nil {
		return nil
	}
	conv := Time(*t)
	return &conv
}
