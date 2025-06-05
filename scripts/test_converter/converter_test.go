package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)





	// Read input file
	inputFile := "random_numbers.txt"
	outputFile := "random_numbers_with_vietnamese.txt"

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		log.Fatalf("Input file %s not found", inputFile)
	}

	// Open input file
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Error opening input file: %v", err)
	}
	defer file.Close()

	// Create output file
	output, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer output.Close()

	startAll := time.Now()

	// Create channels for work distribution
	numbers := make(chan string, 1000)
	results := make(chan TestResult, 1000)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(i+1, numbers, results, &wg)
	}

	// Start result writer goroutine
	var writerWg sync.WaitGroup
	writerWg.Add(1)
	go func() {
		defer writerWg.Done()
		for result := range results {
			if result.Success {
				// Write number and Vietnamese to output file
				if _, err := output.WriteString(fmt.Sprintf("%d %s\n", result.Input, result.Output)); err != nil {
					log.Printf("Error writing result: %v", err)
				}
				log.Printf("Processed: %d (%.2fms)", result.Input, result.DurationMs)
			} else {
				log.Printf("Failed: %d - %s", result.Input, result.Error)
			}
		}
	}()

	// Read numbers and send to workers
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		numberStr := scanner.Text()
		if numberStr == "" {
			continue
		}
		count++
		numbers <- numberStr
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input file: %v", err)
	}

	// Close the numbers channel to signal workers to exit
	close(numbers)

	// Wait for all workers to finish
	wg.Wait()

	// Close results channel after all workers are done
	close(results)

	// Wait for writer to finish
	writerWg.Wait()

	// Wait for writer to finish
	writerWg.Wait()

	totalDuration := time.Since(startAll)
	log.Printf("Test completed. Results written to %s", outputFile)
	log.Printf("Total execution time: %s", totalDuration)
}

func worker(id int, numbers <-chan string, results chan<- TestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for numStr := range numbers {
		// Convert string to int64
		number, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			results <- TestResult{
				Input:   0,
				Success: false,
				Error:   fmt.Sprintf("Invalid number: %v", err),
			}
			continue
		}

		// Create request body
		reqBody := map[string]interface{}{
			"number": number,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			results <- TestResult{
				Input:   number,
				Success: false,
				Error:   fmt.Sprintf("Error creating request: %v", err),
			}
			continue
		}

		start := time.Now()

		// Send request
		resp, err := client.Post(
			baseURL+"/convert",
			"application/json",
			bytes.NewBuffer(jsonBody),
		)

		duration := time.Since(start).Seconds() * 1000 // Convert to milliseconds

		// Handle response
		if err != nil {
			results <- TestResult{
				Input:      number,
				DurationMs: duration,
				Success:    false,
				Error:      fmt.Sprintf("Request failed: %v", err),
			}
			continue
		}
		defer resp.Body.Close()

		// Parse response
		var result struct {
			Number         int64   `json:"number"`
			Vietnamese     string  `json:"vietnamese"`
			ProcessingTimeMs float64 `json:"processing_time_ms"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			results <- TestResult{
				Input:      number,
				DurationMs: duration,
				Success:    false,
				Error:      fmt.Sprintf("Error decoding response: %v", err),
			}
			continue
		}

		// Check if the response is successful
		if resp.StatusCode != http.StatusOK {
			results <- TestResult{
				Input:      number,
				DurationMs: duration,
				Success:    false,
				Error:      fmt.Sprintf("Unexpected status code: %d", resp.StatusCode),
			}
			continue
		}

		// Send successful result
		results <- TestResult{
			Input:      number,
			Output:     result.Vietnamese,
			DurationMs: duration,
			Success:    true,
		}

		// Be nice to the server
		time.Sleep(10 * time.Millisecond)
	}
}
