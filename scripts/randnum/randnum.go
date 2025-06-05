// Package randnum generates random numbers and writes them to a file
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create output file
	file, err := os.Create("random_numbers.txt")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Constants
	const (
		min = 20_000
		max = 20_000_000_000
		count = 100_000
	)

	// Generate and write random numbers
	for i := 0; i < count; i++ {
		// Generate random number between min and max (inclusive)
		randNum := rand.Int63n(max-min+1) + min
		
		// Format: number with newline
		_, err := file.WriteString(strconv.FormatInt(randNum, 10) + "\n")
		if err != nil {
			log.Printf("Error writing to file: %v", err)
			continue
		}

		// Print progress every 10,000 numbers
		if (i+1)%10_000 == 0 {
			fmt.Printf("Generated %d numbers...\n", i+1)
		}
	}

	fmt.Printf("Successfully generated %d random numbers between %d and %d in random_numbers.txt\n", 
		count, min, max)
}
