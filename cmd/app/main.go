package main

import (
	"29th_oct/internal/ds"
	"fmt"
)

func main() {
	fmt.Println("Main function")

	scores := ds.NewSafeMap[string, int]()

	scores.Put("Swavik", 99)

	if score, ok := scores.Get("Swavik"); ok {
		fmt.Println("Swavik's score is: ", score)
	} else {
		fmt.Println("Swavik not found")
	}
}
