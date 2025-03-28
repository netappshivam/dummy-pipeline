package main

import (
	"fmt"
)

func TagHeaderFunc() string {
	//currentYear, currentMonth := time.Now().Year()%100, int(time.Now().Month())
	//currentYearStr := fmt.Sprintf("%d", currentYear)
	//currentMonthStr := fmt.Sprintf("%02d", currentMonth)
	//currentTagHeaderName := currentYearStr + currentMonthStr
	//
	//fmt.Println(currentTagHeaderName)
	//return currentTagHeaderName

	currentTagHeaderName := "v1"

	fmt.Println(currentTagHeaderName)
	return currentTagHeaderName

}
