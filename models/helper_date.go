package models

import (
	"encoding/json"
	"strings"
	"time"

	"4d63.com/tz"
)

const (
	DateLayout = "2006-01-02"
)

type Date time.Time

/*
------------------------
Date Function
------------------------
*/

func NewDateFromString(dateString string) Date {
	loc, _ := tz.LoadLocation("Asia/Bangkok")
	d, err := time.ParseInLocation(DateLayout, dateString, loc)
	if err != nil {
		panic(err)
	}
	return Date(d)
}

func NewDateFromStringWithTime(dateString string) Date {
	loc, _ := tz.LoadLocation("Asia/Bangkok")
	d, err := time.ParseInLocation(TimestampLayout, dateString, loc)
	if err != nil {
		panic(err)
	}
	return Date(d)
}

func NewDateFromTime(t time.Time) Date {
	loc, _ := tz.LoadLocation("Asia/Bangkok")
	d, err := time.ParseInLocation(DateLayout, t.Format(DateLayout), loc)
	if err != nil {
		panic(err)
	}
	return Date(d)
}

func (j *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(DateLayout, s)
	if err != nil {
		return err
	}
	*j = Date(t)
	return nil
}

func (j Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Format(DateLayout))
}

func (j Date) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

func (j Date) Weekday() time.Weekday {
	t := time.Time(j)
	return t.Weekday()
}

func (j Date) String() string {
	return j.Format(DateLayout)
}

func (j *Date) GetBSON() (interface{}, error) {
	if j == nil {
		return nil, nil
	}
	loc, _ := tz.LoadLocation("Asia/Bangkok")
	t := time.Time(*j)
	d, err := time.ParseInLocation(TimestampLayout, t.Format(DateLayout), loc)
	if err != nil {
		return nil, err
	}

	return d, nil
}
