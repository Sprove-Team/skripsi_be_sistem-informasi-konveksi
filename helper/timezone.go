package helper

import (
	"errors"
	"time"

	"github.com/be-sistem-informasi-konveksi/common/message"
)

func GetTimezone(timezone string) (*time.Location,error) {	
	timeLoad, err := time.LoadLocation(timezone)
	if err != nil {
		LogsError(err)
		return nil, errors.New(message.Timezoneunknown)
	}
	return timeLoad, nil
}