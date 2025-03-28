package main

import (
	"fmt"
)

func Xyz() string {
	//currentYear, currentMonth := time.Now().Year()%100, int(time.Now().Month())

	//t := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.UTC)
	//
	//// Calling ISOWeek() method
	//_, week := t.ISOWeek()
	//
	//fmt.Println(week)
	//
	//currentYearStr := fmt.Sprintf("%d", currentYear)
	//currentMonthStr := fmt.Sprintf("%02d", currentMonth)

	//currentTagDateName := currentYearStr + currentMonthStr
	currentTagDateName := "v1"

	fmt.Println(currentTagDateName)
	return currentTagDateName
}
