package compiler

import "time"

type Timestamp struct {
	time.Time
}

var TimestampFormat = "2006-01-02T15:04:05Z"

func NewTimestamp(str string) (Timestamp, error) {
	t, err := time.ParseInLocation(TimestampFormat, str, time.Now().Location())
	if err != nil {
		return Timestamp{}, err
	}

	return Timestamp{t}, nil
}

func (t Timestamp) Equals(u Timestamp) bool {
	return t.Equal(u.Time)
}

func (t Timestamp) GreaterThan(u Timestamp) bool {
	return t.After(u.Time)
}

func (t Timestamp) GreaterThanEquals(u Timestamp) bool {
	if t.Equals(u) {
		return true
	}
	return t.GreaterThan(u)
}

func (t Timestamp) LessThan(u Timestamp) bool {
	return t.Before(u.Time)
}

func (t Timestamp) LessThanEquals(u Timestamp) bool {
	if t.Equals(u) {
		return true
	}
	return t.LessThan(u)
}
