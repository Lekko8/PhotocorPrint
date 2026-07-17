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
var reportMatrix = make(map[ReportKey]SampleData)

// Собирает список файлов в массив имён и запускает чтение данных
func readFileList(filesFolder string) string {
	start := time.Now()

	defer func() {
		if filesList != nil {
			reportMatrix = initDataRead()
		}
	}()

	filesList = nil

	files, err := os.ReadDir(filesFolder) // читаем все файлы в папке
	if err != nil {
		return buildNames(filesList) // если файлов нет, то возвращаем пустой список
	}

	for _, file := range files { // ищем подходящие файлы

		if !file.IsDir() && (strings.ToUpper(filepath.Ext(file.Name())) == ".TXT") {
			filesList = append(filesList, file.Name()) // кладём подходящий
		}
	}
	countOfFiles = len(filesList)

	log.Print("Чтение файлов: ", time.Since(start))
	return buildNames(filesList)
}

// Создаёт .xlsx файл
func createFile() string {
	start := time.Now()

	log.Print("Количество файлов: ", len(filesList))

	_, weekNum := time.Now().ISOWeek()

	resultFileName := "DLS_rlt_" + strconv.Itoa(weekNum) + "_" +
		time.Now().Format("02012006") + "_s_" + ".xlsx"

	xlsx(resultFileName, reportMatrix)

	log.Print("Создание файла: ", time.Since(start))
	return resultFileName
}

// Разбивает массив имён файлов в формат для вывода
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
