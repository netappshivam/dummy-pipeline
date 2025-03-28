package main

import (
	"fmt"
)

func Xyz() string {
	//currentYear, currentMonth := time.Now().Year()%100, int(time.Now().Month())
	//currentYearStr := fmt.Sprintf("%d", currentYear)
	//currentMonthStr := fmt.Sprintf("%02d", currentMonth)
	//currentTagDateName := currentYearStr + currentMonthStr
	//
	//fmt.Println(currentTagDateName)
	//return currentTagDateName

	currentTagDateName := "v1"

	fmt.Println(currentTagDateName)
	return currentTagDateName

}
