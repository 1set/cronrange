package cronrange_test

import (
	"fmt"
	"time"

	"github.com/1set/cronrange"
)

// This example shows greeting according to your local box time.
func Example() {
	crGreetings := make(map[*cronrange.CronRange]string)
	crExprGreetings := map[string]string {
		"DR=360; 0 6 * * *": "Good morning!",
		"DR=360; 0 12 * * *": "Good afternoon!",
		"DR=240; 0 18 * * *": "Good evening!",
		"DR=180; 0 22 * * *": "Good night!",
		"DR=300; 0 1 * * *": "ZzzZZzzzZZZz...",
	}

	for crExpr, greeting := range crExprGreetings {
		if cr, err := cronrange.ParseString(crExpr); err == nil {
			crGreetings[cr] = greeting
		} else {
			fmt.Println("got parse err:", err)
			return
		}
	}

	current := time.Now()
	for cr, greeting := range crGreetings {
		if isWithin, err := cr.IsWithin(current); err == nil {
			if isWithin {
				fmt.Println(greeting)
			}
		} else {
			fmt.Println("got check err:", err)
			return
		}
	}
}
