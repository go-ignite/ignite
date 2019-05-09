package api

import (
	"strconv"
	"time"
)

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	timestamp, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Time(time.Unix(timestamp, 0))
	return nil
}

func (t Time) Time() time.Time {
	return time.Time(t)
}

func (t Time) String() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}
