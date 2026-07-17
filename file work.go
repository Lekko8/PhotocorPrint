package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
)

type WorkerResult struct {
	Key  ReportKey
	Data SampleData
	Err  error
}

// ReportKey информация по файлу
type ReportKey struct {
	Group string // Unique identifier for the order.
	Name  string
}

// SampleData данные из файла
type SampleData struct {
	MeanIntensity     string
	RateOfCorrectUnit string
}

// Запускает горутины чтения .TXT файлов и собирает все данные в WorkerResult
func initDataRead() map[ReportKey]SampleData {

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

	for res := range resultChan {

		if res.Err != nil {
			log.Print(res.Err)
		}

		reportMatrix[res.Key] = res.Data

	}
	return reportMatrix
}

// Собирает данные из .TXT файла во временную мапу
func readFile(fileName string) (ReportKey, SampleData, error) {
	txtFile, err := os.Open(filesFolder + "/" + fileName)
	log.Print("Reading file: ", fileName)
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
	notFound := true
	var mean string
	for scanner.Scan() {
		line := scanner.Text()
		switch notFound {
		case true:
			if strings.Contains(line, "Mean signal intensity") { // 26 строка
				mean = strings.Split(line, ":")[1]
				notFound = false
			}
		case false:
			if strings.Contains(line, "Rate of correct unit") { // 342 строка
				return forKey(fileName), forData(mean, strings.Split(line, ":")[1]), nil
			}
		}

		lineNumber++

	}
	return forKey(fileName), forData("", ""), err
}

func forKey(fileName string) ReportKey {

	names := strings.Split(fileName, "_")

	return ReportKey{names[0], strings.Split(strings.TrimPrefix(fileName, names[0]+"_"), ".")[0]}
}

func forData(mean, line string) SampleData {
	return SampleData{strings.TrimSpace(mean), strings.TrimSpace(line)}
}
