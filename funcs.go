package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func todayMonth() string {
	now := time.Now()
	return monthRu(now.Month())
}

func monthRu(month time.Month) string {
	months := map[time.Month]string{
		time.January:   "Январь",
		time.February:  "Февраль",
		time.March:     "Март",
		time.April:     "Апрель",
		time.May:       "Май",
		time.June:      "Июнь",
		time.July:      "Июль",
		time.August:    "Август",
		time.September: "Сентябрь",
		time.October:   "Октябрь",
		time.November:  "Ноябрь",
		time.December:  "Декабрь",
	}
	return months[month]
}

// Собирает список файлов в массив имён
func readFiles(filesFolder string) string {
	start := time.Now()

	var filesList []string

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

func createXlsx(orderName string) string {
	start := time.Now()

	_, weekNum := time.Now().ISOWeek()

	resultFileName := "DLS_rlt_" + strconv.Itoa(weekNum) + "_" +
		time.Now().Format("02.01.2006") + "_s_" +
		orderName + ".xlsx"

	createFile(resultFileName)
	log.Print("Создание файла: ", time.Since(start))
	return "Файл " + resultFileName + " успешно создан"
}
