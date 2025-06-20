package utils

import (
	"errors"
	"time"
)

func CheckDate(date string) error {
	parsedDate, err := time.Parse("2006-01-02-15-04", date)
	if err != nil {
		return err
	}
	checkDate := parsedDate.UTC().Add(-time.Hour * 3).Unix()
	if checkDate < time.Now().UTC().Add(time.Hour*1).Unix() {
		return errors.New("meet up is out of date")
	}
	return nil
}
