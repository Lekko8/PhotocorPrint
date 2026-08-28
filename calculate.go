package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func calculate(fileName string) string {

	start := time.Now()
	defer func() { log.Print("Обработка файла: ", time.Since(start)) }()

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Print(err.Error())
		return "Ошибка открытия файла " + fileName + ": " + err.Error()
	}
	defer func() {
		if err := f.Save(); err != nil {
			log.Print(err.Error())
		}
		if err := f.Close(); err != nil {
			log.Print(err.Error())
		}

		cmd := exec.Command("cmd", "/c", "start", "", fileName)
		if err := cmd.Start(); err != nil {
			log.Print(err.Error())
		}
	}()

	styles := getStyles(f)

	addHeaderPage(f, styles)

	for _, sheet := range f.GetSheetList() {
		if strings.Contains(sheet, "PBS") ||
			strings.Contains(sheet, "BSA") ||
			strings.Contains(sheet, "latex") {

			singlePeak(f, sheet, styles)
			continue
		}
		if strings.Contains(sheet, "Смесь") {
			doublePeak(f, sheet, styles)
		}
	}

	addResults(f, styles)

	return "Посчитал " + fileName + " за " + time.Since(start).String()
}

func addResults(f *excelize.File, styles map[string]int) {
	const textRStudio = "Расчет проводили с помощью программы R.Studio по скрипту " +
		"Сomparisons_with_control+outliers (cond+ghz+thz+el+nmr)_260529. " +
		"За достоверный результат принимали расчетное значение p.value < 0.05."
	const sheet = "выводы"
	idx, err := f.NewSheet(sheet)
	defer func(idx int, err error) {
		f.SetActiveSheet(idx)
		if err != nil {
			log.Print(err.Error())
		}
	}(idx, err)

	addHeader(f, sheet, styles)

	err = f.SetRowHeight(sheet, 11, 48)
	err = f.SetColWidth(sheet, "A", "A", 15)
	err = f.MergeCell(sheet, "M1", "S4")
	err = f.SetCellStr(sheet, "M1", textRStudio)
	err = f.SetCellStyle(sheet, "M1", "S4", styles["styleRStudio"])
	err = f.MergeCell(sheet, "A11", "A12")
	err = f.MergeCell(sheet, "B11", "B12")
	err = f.MergeCell(sheet, "C11", "C12")
	err = f.MergeCell(sheet, "D11", "D12")
	err = f.MergeCell(sheet, "E11", "E12")
	err = f.MergeCell(sheet, "F11", "J11")
	err = f.MergeCell(sheet, "K11", "M11")
	err = f.MergeCell(sheet, "A13", "A22")
	err = f.MergeCell(sheet, "B13", "B22")
	err = f.SetCellStyle(sheet, "B11", "M11", styles["styleH"])
	err = f.SetCellStyle(sheet, "A11", "M22", styles["style"])

}

func addHeaderPage(f *excelize.File, styles map[string]int) {

	const experiment = "Оценка размеров частиц р-ра"
	const explanation = "Целью эксперимента являлось оценить размеры частиц"
	const sheet = "пробоподготовка"
	err := f.MergeCell(sheet, "A4", "A11")
	err = f.MergeCell(sheet, "B4", "M11")
	err = f.MergeCell(sheet, "B1", "M1")
	err = f.MergeCell(sheet, "B2", "M2")
	err = f.MergeCell(sheet, "B3", "M3")
	err = f.MergeCell(sheet, "B12", "M12")
	err = f.SetCellStr(sheet, "A1", "Эксперимент:")
	err = f.SetCellStr(sheet, "A2", "Дата:")
	err = f.SetCellStr(sheet, "A3", "Исполнитель:")
	err = f.SetCellStr(sheet, "A4", "Описание:")
	err = f.SetCellStr(sheet, "A12", "Пробоподготовка")
	err = f.SetCellStyle(sheet, "A1", "M12", styles["style"])
	err = f.SetColWidth(sheet, "B", "A", 20)
	err = f.SetColWidth(sheet, "B", "M", 17)
	err = f.SetRowHeight(sheet, 1, 39)
	err = f.SetRowHeight(sheet, 3, 23.3)
	err = f.SetRowHeight(sheet, 11, 137.3)
	err = f.SetRowHeight(sheet, 12, 184.5)
	err = f.SetCellStr(sheet, "B1", experiment)
	err = f.SetCellValue(sheet, "B2", time.Now().Format("02.01.2006"))
	err = f.SetCellStr(sheet, "B3", "")
	err = f.SetCellStr(sheet, "B4", explanation)
	err = f.SetCellStr(sheet, "B12", sheet)
	if err != nil {
		log.Print(err.Error())
	}
}

func addHeader(f *excelize.File, sheet string, styles map[string]int) {
	err := f.MergeCell(sheet, "A4", "A8")
	err = f.MergeCell(sheet, "B4", "K8")
	err = f.MergeCell(sheet, "B1", "K1")
	err = f.MergeCell(sheet, "B2", "K2")
	err = f.MergeCell(sheet, "B3", "K3")
	err = f.SetCellFormula(sheet, "A1", "пробоподготовка!A1")
	err = f.SetCellFormula(sheet, "B1", "пробоподготовка!B1")
	err = f.SetCellFormula(sheet, "A2", "пробоподготовка!A2")
	err = f.SetCellFormula(sheet, "B2", "пробоподготовка!B2")
	err = f.SetCellFormula(sheet, "A3", "пробоподготовка!A3")
	err = f.SetCellFormula(sheet, "B3", "пробоподготовка!B3")
	err = f.SetCellFormula(sheet, "A4", "пробоподготовка!A4")
	err = f.SetCellFormula(sheet, "B4", "пробоподготовка!B4")
	err = f.SetCellStyle(sheet, "A1", "K8", styles["style"])
	if err != nil {
		log.Print(err.Error())
	}
}

func singlePeak(f *excelize.File, sheetName string, styles map[string]int) {

	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Print(err.Error())
	}

	var (
		DataTables []map[int]float64
		resultSP   []float64 // R, нм (Position)
	)

	for i, row := range rows {
		if len(row) > 0 && strings.Contains(row[0], "Peak Num") {
			//log.Print(sheetName, "-", i)
			DataTables = append(DataTables, readTable(rows, i))
		}
	}

	for _, table := range DataTables {

		//log.Print("Working on ", table.Area)

		areas := make([]int, 0, len(table))
		positions := make([]float64, 0, len(table))
		for area, position := range table {
			areas = append(areas, area)
			positions = append(positions, position)
		}

		//log.Print("areas: ", areas)

		ans := areas[0]
		for _, area := range areas {
			if area > ans {
				ans = area
			}
		}

		//log.Print("answer: ", ans)

		resultSP = append(resultSP, table[ans])
	}

	//log.Print("resultSP: ", resultSP)

	rateOf := makeRateOf(rows)

	//log.Print("rateOf: ", rateOf)

	addTable(f, sheetName, styles, "G26", 0, resultSP, rateOf)

	addHeader(f, sheetName, styles)

	log.Print("singlePeak done in " + sheetName)
}

func doublePeak(f *excelize.File, sheetName string, styles map[string]int) {

	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Print(err.Error())
	}

	var (
		DataTables  []map[int]float64
		firstPeaks  []float64 // R, нм (Position) для 1 таблицы
		secondPeaks []float64 // R, нм (Position) для 2 таблицы
	)

	for i, row := range rows {
		if len(row) > 0 && strings.Contains(row[0], "Peak Num") {
			//log.Print(sheetName, "-", i)
			DataTables = append(DataTables, readTable(rows, i))
		}
	}

	for _, table := range DataTables {

		//log.Print("Working on ", table.Area)

		var firstPeak float64

		areas := make([]int, 0, len(table))
		positions := make([]float64, 0, len(table))
		for area, position := range table {
			areas = append(areas, area)
			positions = append(positions, position)
		}

		//log.Print("areas: ", areas)

		if len(table) == 1 {
			firstPeaks = append(firstPeaks, positions[0])
			secondPeaks = append(secondPeaks, 0)
			continue
		}

		for j, position := range positions {
			if 3.9 < position && position < 9.0 && areas[j] > 100 {
				firstPeak = position
				break
			}
		}

		maxArea := areas[0]
		maxIdx := 0
		for i, area := range areas {
			if area > maxArea {
				maxArea = area
				maxIdx = i
			}
		}

		if positions[maxIdx] == firstPeak {

			sMaxArea := -1
			secondIdx := -1
			for i, area := range areas {
				if i != maxIdx && area > sMaxArea {
					sMaxArea = area
					secondIdx = i
				}
			}

			if positions[secondIdx] < 30 && areas[secondIdx] > 100 {
				firstPeaks = append(firstPeaks, firstPeak)
				secondPeaks = append(secondPeaks, positions[secondIdx])
				continue
			} else {
				firstPeaks = append(firstPeaks, firstPeak)
				secondPeaks = append(secondPeaks, 0)
				continue
			}

		}

		if positions[maxIdx] < 30 && areas[maxIdx] > 100 {
			firstPeaks = append(firstPeaks, firstPeak)
			secondPeaks = append(secondPeaks, positions[maxIdx])
			continue
		} else {
			firstPeaks = append(firstPeaks, firstPeak)
			secondPeaks = append(secondPeaks, 0)
			continue
		}

	}

	rateOf := makeRateOf(rows)

	//log.Print("rateOf: ", rateOf)

	addTable(f, sheetName, styles, "G26", 1, firstPeaks, rateOf)

	addTable(f, sheetName, styles, "O26", 2, secondPeaks, rateOf)

	addHeader(f, sheetName, styles)

	log.Print("doublePeak done in " + sheetName)
}

func makeRateOf(rows [][]string) []int {

	var rateOf []int
	i := 13
	r := rows[i][4]
	for i < len(rows) {

		//log.Print("row №", i, rows[i])

		rate, err := strconv.Atoi(r)
		if err != nil {
			log.Print(err.Error())
		}

		rateOf = append(rateOf, rate)

		i++

		if len(rows[i]) < 2 {
			break
		}

		r = rows[i][4]
	}
	return rateOf
}

func readTable(dataRows [][]string, i int) map[int]float64 {

	var ans = make(map[int]float64)

	i++
	for i < len(dataRows) {

		//log.Println(dataRows[i])

		if len(dataRows[i]) == 0 {
			break
		}

		if !strings.ContainsAny(dataRows[i][0], "0123456789") {
			break
		}

		//log.Print(dataRows[i][1])

		area, err := strconv.ParseFloat(strings.ReplaceAll(dataRows[i][1], ",", "."), 64)
		if err != nil {
			log.Print(err.Error())
		}
		a := int(area * 1000)

		position, err := strconv.ParseFloat(strings.ReplaceAll(dataRows[i][3], ",", "."), 64)
		if err != nil {
			log.Print(err.Error())
		}

		ans[a] = position

		i++

	}
	//log.Println("Finished makeTable: ", ans)
	return ans
}

func addTable(f *excelize.File, sheetName string, styles map[string]int,
	rCell string, n int, data []float64, rateOf []int) {

	column, row, err := excelize.CellNameToCoordinates(rCell)
	if err != nil {
		log.Print(err.Error())
	}

	l := len(data) - 1

	tlStyleCell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		log.Print(err.Error())
	}
	brStyleCell, err := excelize.CoordinatesToCellName(column+6, row+l+2)
	if err != nil {
		log.Print(err.Error())
	}

	merge(f, sheetName, column, 6, row, 0) // Шапка

	var tableHeadName string

	if n == 0 {
		tableHeadName = "PS"
	} else {
		tableHeadName = "PS пик " + strconv.Itoa(n)
	}

	err = f.SetCellStr(sheetName, rCell, tableHeadName)
	if err != nil {
		log.Print(err.Error())
	}

	row++ // Заполнение заголовков

	tableCollName := []string{"№", "R, нм", "D, нм", "D mean, нм", "SD", "CV", "Rate of correct unit, %"}

	tlhCell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		log.Print(err.Error())
	}

	brhCell, err := excelize.CoordinatesToCellName(column+6, row)
	if err != nil {
		log.Print(err.Error())
	}

	for i, name := range tableCollName {
		cell, err := excelize.CoordinatesToCellName(column+i, row)
		if err != nil {
			log.Print(err.Error())
		}

		err = f.SetCellStr(sheetName, cell, name)
		if err != nil {
			log.Print(err.Error())
		}
	}

	row++ // №

	for i := range data {
		cell, err := excelize.CoordinatesToCellName(column, row+i)
		if err != nil {
			log.Print(err.Error())
		}

		err = f.SetCellInt(sheetName, cell, int64(i+1))
		if err != nil {
			log.Print(err.Error())
		}

		//log.Print(int64(i+1), " in ", cell)
	}

	column++ // R, нм

	for i, d := range data {
		cell, err := excelize.CoordinatesToCellName(column, row+i)
		if err != nil {
			log.Print(err.Error())
		}

		if d == 0 {
			err = f.SetCellValue(sheetName, cell, "-")
			if err != nil {
				log.Print(err.Error())
			}
			continue
		}

		err = f.SetCellValue(sheetName, cell, d)
		if err != nil {
			log.Print(err.Error())
		}

		//log.Print(d, " in ", cell)
	}

	column++ // D, нм

	for i, d := range data {
		cell, err := excelize.CoordinatesToCellName(column, row+i)
		if err != nil {
			log.Print(err.Error())
		}

		dataCell, err := excelize.CoordinatesToCellName(column-1, row+i)
		if err != nil {
			log.Print(err.Error())
		}

		if d == 0 {
			err = f.SetCellValue(sheetName, cell, "-")
			if err != nil {
				log.Print(err.Error())
			}
			continue
		}

		err = f.SetCellFormula(sheetName, cell, dataCell+"*2")
		if err != nil {
			log.Print(err.Error())
		}

		//log.Print(data[i]*2, " in ", cell)
	}

	column++ // D mean, нм

	dCell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		log.Print(err.Error())
	}

	merge(f, sheetName, column, 0, row, l)

	setFormulaGap(f, sheetName, "AVERAGE", column, row, column-1, 0, row, l)

	column++ // SD

	sdCell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		log.Print(err.Error())
	}

	merge(f, sheetName, column, 0, row, l)

	setFormulaGap(f, sheetName, "STDEV", column, row, column-2, 0, row, l)

	column++ // CV

	merge(f, sheetName, column, 0, row, l)

	fCell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		log.Print(err.Error())
	}

	err = f.SetCellFormula(sheetName, fCell, sdCell+"/"+dCell)
	if err != nil {
		log.Print(err.Error())
	}

	column++ // Rate of correct unit, %

	for i, rate := range rateOf {
		cell, err := excelize.CoordinatesToCellName(column, row+i)
		if err != nil {
			log.Print(err.Error())
		}

		err = f.SetCellValue(sheetName, cell, rate)
		if err != nil {
			log.Print(err.Error())
		}
	}

	row += l + 1

	//log.Print("AVERAGE in ", column, row)

	setFormulaGap(f, sheetName, "AVERAGE", column, row, column, 0, row-l-1, l)

	tlRCell, err := excelize.CoordinatesToCellName(column, row-l-1)
	if err != nil {
		log.Print(err.Error())
	}

	brRCell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		log.Print(err.Error())
	}

	log.Printf("%s: %s - %s", sheetName, tlRCell, brRCell)

	err = f.SetCellStyle(sheetName, tlStyleCell, brStyleCell, styles["style"])

	err = f.SetCellStyle(sheetName, tlRCell, brRCell, styles["styleRate"])

	err = f.SetCellStyle(sheetName, fCell, fCell, styles["styleCV"])

	err = f.SetCellStyle(sheetName, tlhCell, brhCell, styles["styleH"])

	err = f.SetCellStyle(sheetName, dCell, sdCell, styles["styleDSD"])
	if err != nil {
		log.Print(err.Error())
	}

	//log.Printf("%s: %s-%s", sheetName, tlStyleCell, brStyleCell)
}

func merge(f *excelize.File, sheetName string, column, cl, row, rl int) {
	mtlCell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		log.Print(err.Error())
	}

	mbrCell, err := excelize.CoordinatesToCellName(column+cl, row+rl)
	if err != nil {
		log.Print(err.Error())
	}

	err = f.MergeCell(sheetName, mtlCell, mbrCell)
	if err != nil {
		log.Print(err.Error())
	}
}

func setFormulaGap(f *excelize.File, sheetName, formula string, column, row, dataColumn, dcl, dataRow, drl int) {
	fCell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		log.Print(err.Error())
	}

	dsCell, err := excelize.CoordinatesToCellName(dataColumn, dataRow)
	if err != nil {
		log.Print(err.Error())
	}

	dfCell, err := excelize.CoordinatesToCellName(dataColumn+dcl, dataRow+drl)
	if err != nil {
		log.Print(err.Error())
	}

	resultFormula := strings.Builder{}
	resultFormula.WriteString(formula)
	resultFormula.WriteString("(")
	resultFormula.WriteString(dsCell)
	resultFormula.WriteString(":")
	resultFormula.WriteString(dfCell)
	resultFormula.WriteString(")")

	err = f.SetCellFormula(sheetName, fCell, resultFormula.String())
	if err != nil {
		log.Print(err.Error())
	}
}

func getStyles(f *excelize.File) map[string]int {

	var (
		styles        = make(map[string]int)
		decimalPlaces = 2
	)
	style, err := f.NewStyle(&excelize.Style{
		DecimalPlaces: &decimalPlaces,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "horizontal", Color: "000000", Style: 1},
			{Type: "vertical", Color: "000000", Style: 1},
		},
	}) // Таблица + округление 2 знака
	styles["style"] = style
	styleH, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",          // Тип заливки
			Color:   []string{"D3D3D3"}, // Светло-серый цвет
			Pattern: 1,                  // Сплошная заливка
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "horizontal", Color: "000000", Style: 1},
			{Type: "vertical", Color: "000000", Style: 1},
		},
	}) // Шапка таблицы
	styles["styleH"] = styleH
	styleCV, err := f.NewStyle(&excelize.Style{
		NumFmt: 9,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "horizontal", Color: "000000", Style: 1},
			{Type: "vertical", Color: "000000", Style: 1},
		},
	}) // CV
	styles["styleCV"] = styleCV
	styleRate, err := f.NewStyle(&excelize.Style{
		NumFmt: 1,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "horizontal", Color: "000000", Style: 1},
			{Type: "vertical", Color: "000000", Style: 1},
		},
	}) // Колонка Rate Of...
	styles["styleRate"] = styleRate
	CustomNumFmtSDS := "0.000"
	styleDSD, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &CustomNumFmtSDS,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "horizontal", Color: "000000", Style: 1},
			{Type: "vertical", Color: "000000", Style: 1},
		},
	}) // D mean и SD
	styles["styleDSD"] = styleDSD
	styleRStudio, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFF2CC"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "horizontal", Color: "000000", Style: 1},
			{Type: "vertical", Color: "000000", Style: 1},
		},
	}) // R.Studio
	styles["styleRStudio"] = styleRStudio
	if err != nil {
		log.Print(err.Error())
	}

	return styles
}

//func setFormula(f *excelize.File, sheetName, formula string, column, row int, dataCells []string) {
//	fCell, err := excelize.CoordinatesToCellName(column, row)
//	if err != nil {
//		log.Print(err.Error())
//	}
//
//	resultFormula := formula + "(" + strings.Join(dataCells, ";") + ")"
//
//	err = f.SetCellFormula(sheetName, fCell, resultFormula)
//	if err != nil {
//		log.Print(err.Error())
//	}
//}
