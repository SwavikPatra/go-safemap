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

	// deleteStatus := scores.Delete("swavik")
	// if deleteStatus {
	// 	fmt.Println("key deleted")
	// } else {
	// 	fmt.Println("Error while deleting key")
	// }

	// scores.Clear()

	// if score, ok := scores.Get("Swavik"); ok {
	// 	fmt.Println("Swavik's score is: ", score)
	// } else {
	// 	fmt.Println("Swavik not found")
	// }

	fmt.Println("Keys: ", scores.Keys())
	scores.Put("Kunu", 98)
	scores.Put("Sonu", 97)

	fmt.Println("Keys: ", scores.Keys())
	fmt.Println("Vals: ", scores.Vals())
}
