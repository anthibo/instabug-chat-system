package helpers

import (
	"fmt"
	"log"
	"time"
)

func Retry[T any](fn func() (T, error), maxRetries int, delay time.Duration) (T, error) {
	var err error
	var res T
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("Retry %d...\n", i+1)
		res, err = fn()
		if err == nil {
			return res, nil
		}
		log.Printf("Retrying after error: %v in %d", err, delay)
		time.Sleep(delay)
	}
	fmt.Println("Operation failed after max retries")
	return res, err
}
