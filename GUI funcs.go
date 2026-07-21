package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var filesList []string
var reportMatrix = make(map[string]SampleData)
var groupsList []string

// Возвращает список файлов в массив имён и запускает чтение данных
func readFileList(filesFolder string) string {
	start := time.Now()

	defer func() {
		if filesList != nil {
			groupsList = groupsList[:0] // очищаем память (сохраняя ёмкость) для чистой перезаписи
			reportMatrix, groupsList = initDataRead()
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
	defer func() { log.Print("Создание файла: ", time.Since(start)) }()

	_, weekNum := time.Now().ISOWeek()

	resultFileName := "DLS_rlt_" + strconv.Itoa(weekNum) + "_" +
		time.Now().Format("02012006") + "_s_" + ".xlsx"

	return xlsx(resultFileName, reportMatrix) + "\nСобрано за " + fmt.Sprint(time.Since(start))
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
