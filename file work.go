package main

import (
	"bufio"
	"log"
	"os"
	"sync"
)

type WorkerResult struct {
	Key  ReportKey
	Data SampleData
	Err  error
}

type ReportKey struct {
	Group string
	Name  string
}

type SampleData struct {
	MeanIntensity     string
	RateOfCorrectUnit string
}

// Запускает горутины чтения .TXT файлов и собирает все данные в WorkerResult
func initDataRead() {

	resultChan := make(chan WorkerResult, len(filesList))

	var wg sync.WaitGroup

	for _, file := range filesList {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			Key, Data, Err := readFile(file)

			resultChan <- WorkerResult{Key: Key, Data: Data, Err: Err}

		}(file)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()
}

// Собирает данные из .TXT файла во временную мапу
func readFile(file string) {
	txtFile, err := os.Open(filesFolder + file)
	if err != nil {
		log.Printf("Error opening file: %v", err)
	}

	defer func() {
		if err := txtFile.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
	scanner := bufio.NewScanner(txtFile)

	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()

		lineNumber++
	}
}
