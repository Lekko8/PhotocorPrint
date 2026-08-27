package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	filesList, groupsList []string
	reportMatrix          = make(map[string]SampleData)
)

// Возвращает список файлов в массив имён и запускает чтение данных
func readFileList(filesFolder string, countOfFiles int) string {
	start := time.Now()

	defer func() {
		if filesList != nil {
			groupsList = groupsList[:0] // очищаем память (сохраняя ёмкость) для чистой перезаписи
			reportMatrix, groupsList = initDataRead(filesFolder)
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
func createFile(filesFolder, order, name string) string {
	start := time.Now()
	defer func() { log.Print("Создание файла: ", time.Since(start)) }()

	return xlsx(filesFolder, createFileName(order, name), reportMatrix) + "\r\nСобрано за " + fmt.Sprint(time.Since(start))
}

func createFileName(order, name string) string {
	_, weekNum := time.Now().ISOWeek()
	ans := strings.Builder{}
	ans.WriteString("DLS_rlt_")
	ans.WriteString(strconv.Itoa(weekNum))
	ans.WriteString("_")
	ans.WriteString(time.Now().Format("060102"))
	if order != "" {
		ans.WriteString("_")
		ans.WriteString(order)
	}
	ans.WriteString("_s_")
	ans.WriteString(name)
	ans.WriteString(".xlsx")
	return ans.String()
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

func getUserName() string {
	corUse, err := user.Current()
	if err != nil {
		log.Print(err.Error())
		return ""
	}
	return strings.Split(corUse.Username, "\\")[1]
}

// Вычисляет сегодняшний месяц
func todayMonth() string {
	now := time.Now()
	return monthRu(now.Month())
}

// Переводит месяц на Ru
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
