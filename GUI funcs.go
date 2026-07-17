package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var filesList []string

// Собирает список файлов в массив имён и запускает чтение данных
func readFileList(filesFolder string) string {
	start := time.Now()

	defer func() {
		if filesList != nil {
			initDataRead()
		}
	}()

	filesList = nil

	files, err := os.ReadDir(filesFolder) // читаем все файлы в папке
	if err != nil {
		countOfFiles = 0
		return buildNames(filesList) // если файлов нет, то возвращаем пустой список
	}

	for _, file := range files { // ищем подходящие файлы

		if !file.IsDir() && (filepath.Ext(file.Name()) == ".TXT" || filepath.Ext(file.Name()) == ".txt") {
			filesList = append(filesList, file.Name()) // кладём подходящий
			countOfFiles++
		}
	}

	log.Print("Чтение файлов: ", time.Since(start))
	return buildNames(filesList)
}

// Создаёт .xlsx файл
func createFile(fileList []string) string {
	start := time.Now()

	log.Print(len(filesList))

	_, weekNum := time.Now().ISOWeek()

	resultFileName := "DLS_rlt_" + strconv.Itoa(weekNum) + "_" +
		time.Now().Format("02012006") + "_s_" + ".xlsx"

	filesWithData := map[string][]string{}

	for _, file := range fileList {
		filesWithData[file] = []string{}
	}

	xlsx(resultFileName)

	log.Print("Создание файла: ", time.Since(start))
	return resultFileName
}

// Расбивает массив имён файлов в формат для вывода
func buildNames(files []string) string {
	if len(files) == 0 {
		return ""
	}
	var filenames strings.Builder
	for _, file := range files {
		if filenames.Len() != 0 {
			filenames.WriteString("\r\n")
		}
		filenames.WriteString(file)
	}
	return filenames.String()
}
