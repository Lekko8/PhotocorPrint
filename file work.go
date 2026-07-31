package main

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type WorkerResult struct {
	FileName string
	Data     SampleData
	Err      error
}

// SampleData данные из файла
type SampleData struct {
	Group             string
	MeanIntensity     float64
	RateOfCorrectUnit int64
}

// Запускает горутины чтения .TXT файлов и собирает все данные в WorkerResult
func initDataRead() (map[string]SampleData, []string) {

	resultChan := make(chan WorkerResult, len(filesList))

	var wg sync.WaitGroup

	for _, fileName := range filesList {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()

			Data, Err := readFile(fileName)

			resultChan <- WorkerResult{FileName: fileName, Data: Data, Err: Err}

		}(fileName)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	tempGroups := make(map[string]struct{})

	for res := range resultChan {
		if res.Err != nil {
			log.Print(res.Err)
		}

		reportMatrix[res.FileName] = res.Data
		tempGroups[res.Data.Group] = struct{}{}
	}

	for group := range tempGroups {
		groupsList = append(groupsList, group)
	}
	slices.Sort(groupsList)
	log.Println(groupsList)

	return reportMatrix, groupsList
}

// Собирает данные из .TXT файла во временную мапу
func readFile(fileName string) (SampleData, error) {

	txtFile, err := os.Open(filesFolder + "/" + fileName)

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

		if notFound {
			if strings.Contains(line, "Mean signal intensity") { // 26 строка
				mean = strings.Split(line, ":")[1]
				notFound = false
			}
		} else {
			if strings.Contains(line, "Rate of correct unit") { // 342 строка
				return forData(fileName, mean, strings.Split(line, ":")[1]), nil
			}
		}

		lineNumber++
	}
	return forData(fileName, "", ""), err
}

// Обрабатывает прочитанные данные
func forData(fileName, mean, line string) SampleData {

	meanF, err := strconv.ParseFloat(strings.TrimSpace(mean), 64)
	if err != nil {
		log.Print(err)
	}

	rateOf, err := strconv.ParseInt(strings.TrimSuffix(strings.TrimSpace(line), "."), 0, 64)
	if err != nil {
		log.Print(err)
	}

	return SampleData{strings.Split(fileName, "_")[0],
		meanF,
		rateOf}
}
