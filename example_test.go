package cronrange_test

import (
	"fmt"
	"time"

	"github.com/1set/cronrange"
)

// This example shows greeting according to your local box time.
func ExampleCronRange_IsWithin() {
	crGreetings := make(map[*cronrange.CronRange]string)
	crExprGreetings := map[string]string{
		"DR=360; 0 6 * * *":  "Good morning!",
		"DR=360; 0 12 * * *": "Good afternoon!",
		"DR=240; 0 18 * * *": "Good evening!",
		"DR=180; 0 22 * * *": "Good night!",
		"DR=300; 0 1 * * *":  "ZzzZZzzzZZZz...",
	}

	// create cronrange from expressions
	for crExpr, greeting := range crExprGreetings {
		if cr, err := cronrange.ParseString(crExpr); err == nil {
			crGreetings[cr] = greeting
		} else {
			fmt.Println("got parse err:", err)
			return
		}
	}

	// check if current time fails in any time range
	current := time.Now()
	for cr, greeting := range crGreetings {
		if isWithin, err := cr.IsWithin(current); err == nil {
			if isWithin {
				fmt.Println(greeting)
			}
		} else {
			fmt.Println("got check err:", err)
			break
		}
	}
}

// This example lists next 5 daily happy hours of Lava Lava Beach Club after 2019.11.09.
func ExampleCronRange_NextOccurrences() {
	cr, err := cronrange.New("0 15 * * *", "Pacific/Honolulu", 120)
	if err != nil {
		fmt.Println("fail to create:", err)
		return
	}

	loc, _ := time.LoadLocation("Pacific/Honolulu")
	currTime := time.Date(2019, 11, 9, 16, 55, 0, 0, loc)
	happyHours, err := cr.NextOccurrences(currTime, 5)
	for _, happyHour := range happyHours {
		fmt.Println(happyHour)
	}

	// Output:
	// [2019-11-10T15:00:00-10:00,2019-11-10T17:00:00-10:00]
	// [2019-11-11T15:00:00-10:00,2019-11-11T17:00:00-10:00]
	// [2019-11-12T15:00:00-10:00,2019-11-12T17:00:00-10:00]
	// [2019-11-13T15:00:00-10:00,2019-11-13T17:00:00-10:00]
	// [2019-11-14T15:00:00-10:00,2019-11-14T17:00:00-10:00]
}
