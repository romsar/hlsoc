package grpc

import (
	"google.golang.org/genproto/googleapis/type/date"
	"time"
)

func fromProtoDate(d *date.Date) time.Time {
	return time.Date(
		int(d.GetYear()),
		time.Month(d.GetMonth()),
		int(d.GetDay()),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func toProtoDate(d time.Time) *date.Date {
	return &date.Date{
		Year:  int32(d.Year()),
		Month: int32(d.Month()),
		Day:   int32(d.Day()),
	}
}
