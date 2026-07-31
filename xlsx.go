package main

import (
	_ "image/jpeg"
	"log"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Собирает .xlsx файл (с 115 строки частично сгенерировано ИИ)
func xlsx(resultFileName string, reportMatrix map[string]SampleData) string {

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Ошибка при создании файла %v", err)
		}
	}()

	sd, err := excelize.OpenFile(filesFolder + "\\Сырые данные.xlsx")
	if err != nil {
		log.Println("Ошибка открытия сырых данных", err)
		return "Ошибка открытия сырых данных: " + err.Error()
	}
	defer func() {
		if err := sd.Close(); err != nil {
			log.Printf("Ошибка с файлом сырых данных %v", err)
		}
	}()

	picsCellsF := []string{"G10", "L10", "L18"}

	for idx := range groupsList {
		sheetName := groupsList[idx]
		log.Println(sheetName)

		index, err := f.NewSheet(sheetName)
		if err != nil {
			log.Printf("Ошибка создания листа %s: %v", sheetName, err)
			continue
		}

		m, err := sd.GetMergeCells(sheetName, true)
		if err != nil {
			log.Print("GetMergeCells:", err.Error())
			return err.Error()
		}
		if len(m) < 2 {
			log.Printf("Недостаточно объединенных ячеек для %s", sheetName)
			continue
		}

		p2 := m[0].GetStartAxis()
		p1 := m[1].GetStartAxis()
		col, row, err := excelize.CellNameToCoordinates(p2)
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
		p3, err := excelize.CoordinatesToCellName(col, row+1)
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}

		picsCellSD := []string{p1, p2, p3}

		var pics []excelize.Picture
		for _, c := range picsCellSD {
			pic, err := sd.GetPictures(sheetName, c)
			if err != nil {
				log.Print("Ошибка получения картинки", err)
				continue
			}
			if len(pic) > 0 {
				pics = append(pics, pic[0])
			}
		}

		for i, pic := range pics {
			if i < len(picsCellsF) {
				err = f.AddPictureFromBytes(sheetName, picsCellsF[i], &excelize.Picture{
					Extension: pic.Extension,
					File:      pic.File,
					Format: &excelize.GraphicOptions{
						ScaleX:          0.7,
						ScaleY:          0.7,
						LockAspectRatio: true,
					},
				})
				if err != nil {
					log.Print("Ошибка вставки картинки: ", err)
					return "Ошибка вставки картинки: " + err.Error()
				}
			}
		}

		rows, err := sd.GetRows(sheetName)
		if err != nil {
			log.Print(err)
			continue
		}

		_, startReadingRow, err := excelize.CellNameToCoordinates(p3)
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
		startReadingRow = startReadingRow + 3

		currentRow := 10

		// ============ 1. НАХОДИМ ПЕРВЫЙ МАРКЕР File: ============
		firstFileRow := 0
		for rowNum := 1; rowNum < len(rows); rowNum++ {
			if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "File: ") {
				firstFileRow = rowNum
				break
			}
		}

		if firstFileRow == 0 {
			f.SetActiveSheet(index)
			continue
		}

		// ============ 2. КОПИРУЕМ ШАПКУ ============
		var prevVal string
		for rowNum := 1; rowNum < firstFileRow && rowNum < len(rows); rowNum++ {
			srcCell, _ := excelize.CoordinatesToCellName(1, rowNum)
			val, _ := sd.GetCellValue(sheetName, srcCell)

			if val == "" {
				if strings.Contains(prevVal, "File(s):") {
					richText, _ := sd.GetCellRichText(sheetName, srcCell)
					dstCell, _ := excelize.CoordinatesToCellName(1, currentRow)
					if len(richText) > 0 {
						_ = f.SetCellRichText(sheetName, dstCell, richText)
					} else {
						_ = f.SetCellValue(sheetName, dstCell, val)
					}
					copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
					currentRow++
				}
				prevVal = val
				continue
			}

			richText, _ := sd.GetCellRichText(sheetName, srcCell)
			dstCell, _ := excelize.CoordinatesToCellName(1, currentRow)

			if len(richText) > 0 {
				_ = f.SetCellRichText(sheetName, dstCell, richText)
			} else {
				_ = f.SetCellValue(sheetName, dstCell, val)
			}
			copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
			currentRow++
			prevVal = val
		}

		// Пустая строка после шапки
		currentRow++

		// ============ 3. НАХОДИМ ЗАГОЛОВКИ ============
		var headerRows []int
		for rowNum := firstFileRow; rowNum < len(rows); rowNum++ {
			if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "Distribution analysis") {
				for j := rowNum - 1; j <= rowNum+5 && j < len(rows); j++ {
					if j >= firstFileRow && j < len(rows) {
						if len(rows[j]) > 0 && !strings.Contains(rows[j][0], "File: ") && rows[j][0] != "" {
							headerRows = append(headerRows, j)
						}
					}
				}
				break
			}
		}

		// ============ 4. ЗАПИСЫВАЕМ ЗАГОЛОВКИ ============
		for _, headerRow := range headerRows {
			if headerRow < len(rows) && len(rows[headerRow]) > 0 {
				for colIdx := 0; colIdx < 5 && colIdx < len(rows[headerRow]); colIdx++ {
					if rows[headerRow][colIdx] != "" {
						srcCell, _ := excelize.CoordinatesToCellName(colIdx+1, headerRow+1)
						dstCell, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)

						richText, _ := sd.GetCellRichText(sheetName, srcCell)
						if len(richText) > 0 {
							_ = f.SetCellRichText(sheetName, dstCell, richText)
						} else {
							_ = f.SetCellValue(sheetName, dstCell, rows[headerRow][colIdx])
						}
						copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
					}
				}
				currentRow++
			}
		}
		currentRow++ // Пустая строка после заголовков

		// ============ 5. НАХОДИМ ВСЕ БЛОКИ ============
		type Block struct {
			markerRow  int
			tableStart int
			tableEnd   int
			c2Row      int
		}
		var blocks []Block

		for rowNum := firstFileRow; rowNum < len(rows); rowNum++ {
			if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "File: ") {
				block := Block{markerRow: rowNum}

				for j := rowNum + 1; j < len(rows); j++ {
					if len(rows[j]) > 0 && strings.Contains(rows[j][0], "Peak Num") {
						if j+1 < len(rows) && len(rows[j+1]) > 0 && len(rows[j+1]) >= 2 {
							hasData := false
							for colIdx := 1; colIdx < len(rows[j+1]) && colIdx < 5; colIdx++ {
								if rows[j+1][colIdx] != "" {
									if _, err := strconv.ParseFloat(strings.Replace(rows[j+1][colIdx], ",", ".", -1), 64); err == nil {
										hasData = true
										break
									}
								}
							}
							if hasData {
								block.tableStart = j
								break
							}
						}
					}
				}

				if block.tableStart == 0 {
					continue
				}

				for j := block.tableStart + 1; j < len(rows); j++ {
					if len(rows[j]) == 0 {
						continue
					}
					if strings.HasPrefix(rows[j][0], "c") {
						block.c2Row = j
						block.tableEnd = j - 1
						break
					}
				}

				if block.tableEnd == 0 {
					block.tableEnd = block.tableStart + 5
				}

				blocks = append(blocks, block)
			}
		}

		// ============ 6. ЗАПИСЫВАЕМ ВСЕ БЛОКИ ============
		for _, block := range blocks {
			// Маркер
			if block.markerRow < len(rows) && len(rows[block.markerRow]) > 0 {
				srcCell, _ := excelize.CoordinatesToCellName(1, block.markerRow+1)
				dstCell, _ := excelize.CoordinatesToCellName(1, currentRow)
				richText, _ := sd.GetCellRichText(sheetName, srcCell)
				if len(richText) > 0 {
					_ = f.SetCellRichText(sheetName, dstCell, richText)
				} else {
					_ = f.SetCellValue(sheetName, dstCell, rows[block.markerRow][0])
				}
				copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
				currentRow++
			}

			// Пустая строка после маркера
			currentRow++

			// Таблица
			for rowIdx := block.tableStart; rowIdx <= block.tableEnd && rowIdx < len(rows); rowIdx++ {
				if len(rows[rowIdx]) == 0 {
					continue
				}
				for colIdx := 0; colIdx < len(rows[rowIdx]) && colIdx < 5; colIdx++ {
					srcCell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
					dstCell, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)

					val := rows[rowIdx][colIdx]
					if colIdx > 0 && val != "" {
						if _, err := strconv.ParseFloat(strings.Replace(val, ",", ".", -1), 64); err == nil {
							val = strings.Replace(val, ".", ",", -1)
						}
					}

					richText, _ := sd.GetCellRichText(sheetName, srcCell)
					if len(richText) > 0 && strings.Contains(rows[rowIdx][colIdx], "c") {
						_ = f.SetCellRichText(sheetName, dstCell, richText)
					} else {
						_ = f.SetCellValue(sheetName, dstCell, val)
					}
					copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
				}
				currentRow++
			}

			// x2: в колонку F
			if block.c2Row > 0 && block.c2Row < len(rows) && len(rows[block.c2Row]) > 0 {
				c2Row := currentRow - 1
				srcCell, _ := excelize.CoordinatesToCellName(1, block.c2Row+1)
				dstCell, _ := excelize.CoordinatesToCellName(6, c2Row)

				richText, _ := sd.GetCellRichText(sheetName, srcCell)
				if len(richText) > 0 {
					_ = f.SetCellRichText(sheetName, dstCell, richText)
				} else {
					_ = f.SetCellValue(sheetName, dstCell, rows[block.c2Row][0])
				}
				copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
			}

			// Пустая строка после блока
			currentRow++
		}

		// ============ 7. ДАННЫЕ ИЗ REPORT MATRIX ============
		_ = f.SetCellValue(sheetName, "D13", "Mean intensity")
		_ = f.SetCellValue(sheetName, "E13", "Rate of correct unit")

		for i := 0; true; i++ {
			cellName := "A" + strconv.Itoa(i+14)
			contCell, err := f.GetCellValue(sheetName, cellName)
			if err != nil || contCell == "" {
				break
			}
			contCell = strings.TrimPrefix(contCell, "¨ ")
			if data, exists := reportMatrix[contCell]; exists {
				_ = f.SetCellValue(sheetName, "D"+strconv.Itoa(i+14), data.MeanIntensity)
				_ = f.SetCellValue(sheetName, "E"+strconv.Itoa(i+14), data.RateOfCorrectUnit)
			}
		}

		_ = f.SetColWidth(sheetName, "A", "A", 23.14)
		_ = f.SetColWidth(sheetName, "D", "F", 15)
		f.SetActiveSheet(index)
	}

	_ = f.SetSheetName("Sheet1", "пробоподготовка")
	_ = f.SaveAs(resultFileName)

	return "Файл " + resultFileName + " успешно создан"
}

func copyCellStyle(src, dst *excelize.File, srcSheet, dstSheet, srcCell, dstCell string) {
	if srcCell == "" || dstCell == "" {
		return
	}

	styleID, err := src.GetCellStyle(srcSheet, srcCell)
	if err != nil || styleID == 0 {
		return
	}

	style, err := src.GetStyle(styleID)
	if err != nil || style == nil {
		return
	}

	newStyleID, err := dst.NewStyle(style)
	if err != nil {
		return
	}

	_ = dst.SetCellStyle(dstSheet, dstCell, dstCell, newStyleID)
}
