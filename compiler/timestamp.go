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
