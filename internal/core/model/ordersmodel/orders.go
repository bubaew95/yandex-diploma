package ordersmodel

import (
	"fmt"
	"time"
)

const TimeFormat = "02-01-2006 15:04:05"

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", ct.Format(TimeFormat))
	return []byte(formatted), nil
}

func (ct CustomTime) UnmarshalJSON(b []byte) (err error) {
	parsedTime, err := time.Parse(`"`+TimeFormat+`"`, string(b))
	if err != nil {
		return err
	}
	ct.Time = parsedTime

	return nil
}

type Orders struct {
	Number     int64      `json:"number"`
	Status     string     `json:"status"`
	Accrual    int        `json:"accrual"`
	UploadedAt CustomTime `json:"uploaded_at"`
}
