package main

import (
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

func search(filesFolder string) string {
	return "Прочитаны файлы из " + filesFolder
}

func createXlsx(orderName string) string {
	_, weekNum := time.Now().ISOWeek()
	//orderName := inputOrderName.Text()
	//files := dataCreate(capturedFiles)
	//results := calculateFiles(files)
	//createrXlsx(files, results, orderName, capturedDate)
	resultFileName := "DLS_rlt_" + strconv.Itoa(weekNum) + "_" + time.Now().Format("02.01.2006") + "_s_" + orderName + ".xlsx"

	return "Файл " + resultFileName + " успешно создан"
}

// читает файлы и собирает из них данные (кроме измерений)
func readFiles(filesFolder string) {

	/*var sampleFiles []string
	//var rawTxtFiles []os.DirEntry

	files, err := os.ReadDir(filesFolder) // читаем все файлы в папке
	if err != nil {
		return sampleFiles // если файлов нет, то возвращаем пустой список
	}
	for _, file := range files { // ищем подходящие файлы
		if !file.IsDir() && filepath.Ext(file.Name()) == ".csv" && slices.Contains(strings.Split(targetDate, " "), fileDateTime(file.Name()).Format("02.01.2006")) {
			filteredFiles = append(filteredFiles, file) // кладём подходящий
		}
	}
	sorting(filteredFiles) // сортируем
	for _, file := range filteredFiles {
		sampleFiles = append(sampleFiles, createFileWData(file)) // заполняем список файлов файлами
	}

	return sampleFiles // []FileWData без измеренных данных*/
}
