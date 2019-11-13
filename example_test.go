package cronrange_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/1set/cronrange"
)

// This example creates an instance representing every New Year's Day in Tokyo.
func ExampleNew() {
	cr, err := cronrange.New("0 0 1 1 *", "Asia/Tokyo", 60*24)
	if err != nil {
		fmt.Println("fail to create:", err)
		return
	}

	fmt.Println(cr)
	// Output: DR=1440; TZ=Asia/Tokyo; 0 0 1 1 *
}

// This example creates an instance with an expression representing every New Year's Day in Tokyo.
func ExampleParseString() {
	cr, err := cronrange.ParseString("DR=1440;TZ=Asia/Tokyo;0 0 1 1 *")
	if err != nil {
		fmt.Println("fail to create:", err)
		return
	}

	fmt.Println(cr)
	// Output: DR=1440; TZ=Asia/Tokyo; 0 0 1 1 *
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

// This example demonstrates serializing a struct containing CronRange to JSON.
func ExampleCronRange_MarshalJSON() {
	cr, err := cronrange.ParseString("DR=240;TZ=America/New_York;0 8 1 1 *")
	if err != nil {
		fmt.Println("got parse err:", err)
		return
	}
	ss := struct {
		Expr *cronrange.CronRange
		Num  int
		Goal string
	}{cr, 42, "Morning"}

	if bytes, err := json.Marshal(ss); err == nil {
		fmt.Println(string(bytes))
	} else {
		fmt.Println("got marshal err:", err)
	}

	// Output: {"Expr":"DR=240; TZ=America/New_York; 0 8 1 1 *","Num":42,"Goal":"Morning"}
}
