package helper

import (
	"time"
	_ "time/tzdata"
)

// func GetTimezone(timezone string) (*time.Location, error) {
// 	timeLoad, err := time.LoadLocation(timezone)
// 	if err != nil {
// 		LogsError(err)
// 		return nil, errors.New(message.Timezoneunknown)
// 	}
// 	return timeLoad, nil
// }

func GetStartEndUTC(start, end string) (*time.Time, *time.Time, error) {
	endDate, err := time.Parse(time.DateOnly, end)
	if err != nil {
		LogsError(err)
		return nil, nil, err
	}
	startDate, err := time.Parse(time.DateOnly, start)
	if err != nil {
		LogsError(err)
		return nil, nil, err
	}

	startDateUTC := startDate.UTC()
	endDateUTC := endDate.UTC()

	return &startDateUTC, &endDateUTC, nil
}
