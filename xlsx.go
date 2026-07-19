package main

import (
	_ "image/jpeg"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

func xlsx(resultFileName string, reportMatrix map[string]SampleData) {
	log.Println(reportMatrix)
	log.Println(groupsList)

	// Создаем новый файл
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Ошибка при создании файла %v", err)
		}
	}()

	// Открываем исходный файл
	sd, err := excelize.OpenFile(filesFolder + "\\Сырые данные.xlsx")
	if err != nil {
		log.Panic("Ошибка открытия сырых данных", err)
	}
	defer func() {
		if err := sd.Close(); err != nil {
			log.Printf("Ошибка с файлом сырых данных %v", err)
		}
	}()

	// Ячейки для вставки картинок в новом файле
	picsCellsF := []string{"G13", "N13", "N24"}

	for idx := range groupsList {
		sheetName := groupsList[idx]

		// Создаем новый лист
		index, err := f.NewSheet(sheetName)
		if err != nil {
			log.Printf("Ошибка создания листа %s: %v", sheetName, err)
			continue
		}

		// Получаем объединенные ячейки для определения структуры
		m, err := sd.GetMergeCells(sheetName, true)
		if len(m) < 2 {
			log.Printf("Недостаточно объединенных ячеек для %s", sheetName)
			continue
		}

		p2 := m[0].GetStartAxis()
		p1 := m[1].GetStartAxis()
		col, row, _ := excelize.CellNameToCoordinates(p2)
		p3, _ := excelize.CoordinatesToCellName(col, row+1)

		// ============ КОПИРУЕМ КАРТИНКИ ============
		picsCellSD := []string{p1, p2, p3}
		log.Println(sheetName, "Ячейки чтения картинок: ", picsCellSD)

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
					Format:    pic.Format,
				})
				if err != nil {
					log.Print("Ошибка вставки картинки: ", err)
				}
			}
		}

		// ============ 1. СНАЧАЛА ВСТАВЛЯЕМ 6 ПУСТЫХ СТРОК ============
		//err = f.InsertRows(sheetName, 1, 6)
		//if err != nil {
		//	log.Printf("Ошибка вставки строк: %v", err)
		//}

		// ============ ПОЛУЧАЕМ ВСЕ СТРОКИ ИЗ ИСХОДНОГО ФАЙЛА ============
		rows, err := sd.GetRows(sheetName)
		if err != nil {
			log.Print(err)
			continue
		}

		_, startReadingRow, _ := excelize.CellNameToCoordinates(p3)
		startReadingRow = startReadingRow + 3

		// ============ 2. КОПИРУЕМ ШАПКУ И МАРКЕРЫ ============
		currentRow := 10

		for rowNum := 1; rowNum < startReadingRow; rowNum++ {
			srcCell, _ := excelize.CoordinatesToCellName(1, rowNum)
			val, err := sd.GetCellValue(sheetName, srcCell)
			if err != nil {
				continue
			}

			if val == "" {
				continue
			}

			richText, _ := sd.GetCellRichText(sheetName, srcCell)
			dstCell, _ := excelize.CoordinatesToCellName(1, currentRow)

			if len(richText) > 0 {
				f.SetCellRichText(sheetName, dstCell, richText)
			} else {
				f.SetCellValue(sheetName, dstCell, val)
			}

			copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
			currentRow++
		}

		// Добавляем пустую строку после шапки (после "File(s): ...")
		rowsAfterCopy, _ := f.GetRows(sheetName)
		for i, row := range rowsAfterCopy {
			if len(row) > 0 && strings.Contains(row[0], "File(s):") {
				f.InsertRows(sheetName, i+2, 1)
				break
			}
		}

		// Запоминаем строку, где закончились маркеры
		headerEndRow := currentRow - 1

		if headerEndRow < 19 {
			headerEndRow = currentRow - 1
		}

		log.Printf("%s: Последний маркер в строке %d", sheetName, headerEndRow)

		// ============ 3. КОПИРУЕМ ОСНОВНЫЕ ДАННЫЕ ============
		// Находим блок заголовков
		var headerRows []int
		for rowNum := startReadingRow; rowNum < len(rows); rowNum++ {
			if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "Distribution analysis") {
				for j := rowNum - 1; j <= rowNum+5 && j < len(rows); j++ {
					if j >= startReadingRow {
						headerRows = append(headerRows, j)
					}
				}
				break
			}
		}

		// Находим все блоки с данными
		type Block struct {
			markerRow  int
			tableStart int
			tableEnd   int
			c2Row      int
		}
		var blocks []Block

		for rowNum := startReadingRow; rowNum < len(rows); rowNum++ {
			if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "File: ") {
				block := Block{markerRow: rowNum}

				for j := rowNum + 2; j < len(rows); j++ {
					if len(rows[j]) > 0 && strings.Contains(rows[j][0], "Peak Num") {
						block.tableStart = j
						break
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

		// ============ 4. ВСТАВЛЯЕМ ДАННЫЕ ПОСЛЕ МАРКЕРОВ ============
		dataStartRow := headerEndRow + 2
		currentRow = dataStartRow

		if len(blocks) > 0 {
			// Записываем заголовки через RichText
			for _, headerRow := range headerRows {
				if headerRow < len(rows) {
					for colIdx := 0; colIdx < 5 && colIdx < len(rows[headerRow]); colIdx++ {
						if rows[headerRow][colIdx] != "" {
							srcCell, _ := excelize.CoordinatesToCellName(colIdx+1, headerRow+1)
							dstCell, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)

							richText, err := sd.GetCellRichText(sheetName, srcCell)
							if err == nil && len(richText) > 0 {
								f.SetCellRichText(sheetName, dstCell, richText)
							} else {
								f.SetCellValue(sheetName, dstCell, rows[headerRow][colIdx])
							}
							copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
						}
					}
					currentRow++
				}
			}

			// Пустая строка после заголовков
			currentRow++

			// Записываем блоки
			for _, block := range blocks {
				// Маркер
				if block.markerRow < len(rows) {
					srcCell, _ := excelize.CoordinatesToCellName(1, block.markerRow+1)
					dstCell, _ := excelize.CoordinatesToCellName(1, currentRow)

					richText, err := sd.GetCellRichText(sheetName, srcCell)
					if err == nil && len(richText) > 0 {
						f.SetCellRichText(sheetName, dstCell, richText)
					} else {
						f.SetCellValue(sheetName, dstCell, rows[block.markerRow][0])
					}
					copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
					currentRow++
				}

				// Пустая строка
				currentRow++

				// Таблица
				for rowIdx := block.tableStart; rowIdx <= block.tableEnd && rowIdx < len(rows); rowIdx++ {
					if len(rows[rowIdx]) > 0 {
						for colIdx := 0; colIdx < len(rows[rowIdx]) && colIdx < 5; colIdx++ {
							if rows[rowIdx][colIdx] != "" {
								srcCell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
								dstCell, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)

								richText, err := sd.GetCellRichText(sheetName, srcCell)
								if err == nil && len(richText) > 0 && strings.Contains(rows[rowIdx][colIdx], "c") {
									f.SetCellRichText(sheetName, dstCell, richText)
								} else {
									f.SetCellValue(sheetName, dstCell, rows[rowIdx][colIdx])
								}

								copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
							}
						}
						currentRow++
					}
				}

				// c 2: в колонку F - через RichText
				if block.c2Row > 0 && block.c2Row < len(rows) && len(rows[block.c2Row]) > 0 {
					c2Row := currentRow - 1
					srcCell, _ := excelize.CoordinatesToCellName(1, block.c2Row+1)
					dstCell, _ := excelize.CoordinatesToCellName(6, c2Row)

					richText, err := sd.GetCellRichText(sheetName, srcCell)
					if err == nil && len(richText) > 0 {
						f.SetCellRichText(sheetName, dstCell, richText)
					} else {
						f.SetCellValue(sheetName, dstCell, rows[block.c2Row][0])
					}
					copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
				}

				// Пустая строка после блока
				currentRow++
			}
		} else {
			// ============ ЕСЛИ НЕТ БЛОКОВ (latex с 1 файлом) ============
			log.Printf("%s: Нет блоков, копируем данные напрямую", sheetName)

			// Находим маркер File:
			var fileMarker string
			var fileMarkerRow int
			for rowNum := startReadingRow; rowNum < len(rows); rowNum++ {
				if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "File: ") {
					fileMarker = rows[rowNum][0]
					fileMarkerRow = rowNum
					log.Printf("%s: Найден маркер File: в строке %d: %s", sheetName, rowNum, fileMarker)
					break
				}
			}

			// Если маркер не найден, ищем в строках до startReadingRow
			if fileMarker == "" {
				for rowNum := 0; rowNum < startReadingRow; rowNum++ {
					if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "File: ") {
						fileMarker = rows[rowNum][0]
						fileMarkerRow = rowNum
						log.Printf("%s: Найден маркер File: в строке %d (до startReadingRow): %s", sheetName, rowNum, fileMarker)
						break
					}
				}
			}

			// Находим данные для latex (Distribution analysis и окружающие строки)
			var dataRows []int
			for rowNum := startReadingRow; rowNum < len(rows); rowNum++ {
				if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "Distribution analysis") {
					log.Printf("%s: Найдена Distribution analysis в строке %d", sheetName, rowNum)
					for j := rowNum - 1; j <= rowNum+5 && j < len(rows); j++ {
						if j >= startReadingRow {
							dataRows = append(dataRows, j)
							log.Printf("%s: Добавлена строка %d для заголовков", sheetName, j)
						}
					}
					break
				}
			}

			// Ищем таблицу
			var tableStart, tableEnd, c2Row int
			for rowNum := startReadingRow; rowNum < len(rows); rowNum++ {
				if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "Peak Num") {
					tableStart = rowNum
					log.Printf("%s: Найдена таблица в строке %d", sheetName, rowNum)
					for j := tableStart + 1; j < len(rows); j++ {
						if len(rows[j]) == 0 {
							continue
						}
						if strings.HasPrefix(rows[j][0], "c") {
							c2Row = j
							tableEnd = j - 1
							log.Printf("%s: Найдена c 2: в строке %d, конец таблицы %d", sheetName, c2Row, tableEnd)
							break
						}
					}
					break
				}
			}

			if tableStart > 0 && fileMarker != "" {
				log.Printf("%s: Начинаем запись данных для latex", sheetName)

				// 1. Добавляем пустую строку перед маркером (как в основной логике)
				currentRow++

				// 2. Записываем маркер File:
				srcCell, _ := excelize.CoordinatesToCellName(1, fileMarkerRow+1)
				dstCell, _ := excelize.CoordinatesToCellName(1, currentRow)

				richText, err := sd.GetCellRichText(sheetName, srcCell)
				if err == nil && len(richText) > 0 {
					f.SetCellRichText(sheetName, dstCell, richText)
				} else {
					f.SetCellValue(sheetName, dstCell, fileMarker)
				}
				copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
				currentRow++
				log.Printf("%s: Записан маркер в строку %d", sheetName, currentRow-1)

				// 3. Пустая строка после маркера
				currentRow++

				// 4. Записываем заголовки (Distribution analysis...)
				for _, rowIdx := range dataRows {
					if rowIdx < len(rows) {
						for colIdx := 0; colIdx < 5 && colIdx < len(rows[rowIdx]); colIdx++ {
							if rows[rowIdx][colIdx] != "" {
								srcCell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
								dstCell, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)

								richText, err := sd.GetCellRichText(sheetName, srcCell)
								if err == nil && len(richText) > 0 {
									f.SetCellRichText(sheetName, dstCell, richText)
								} else {
									f.SetCellValue(sheetName, dstCell, rows[rowIdx][colIdx])
								}
								copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
							}
						}
						currentRow++
					}
				}

				// 5. Пустая строка после заголовков
				currentRow++

				// 6. Записываем таблицу
				for rowIdx := tableStart; rowIdx <= tableEnd && rowIdx < len(rows); rowIdx++ {
					if len(rows[rowIdx]) > 0 {
						for colIdx := 0; colIdx < len(rows[rowIdx]) && colIdx < 5; colIdx++ {
							if rows[rowIdx][colIdx] != "" {
								srcCell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
								dstCell, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)

								richText, err := sd.GetCellRichText(sheetName, srcCell)
								if err == nil && len(richText) > 0 && strings.Contains(rows[rowIdx][colIdx], "c") {
									f.SetCellRichText(sheetName, dstCell, richText)
								} else {
									f.SetCellValue(sheetName, dstCell, rows[rowIdx][colIdx])
								}
								copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
							}
						}
						currentRow++
					}
				}

				// 7. c 2: в колонку F
				if c2Row > 0 && c2Row < len(rows) && len(rows[c2Row]) > 0 {
					c2RowDst := currentRow - 1
					srcCell, _ := excelize.CoordinatesToCellName(1, c2Row+1)
					dstCell, _ := excelize.CoordinatesToCellName(6, c2RowDst)

					richText, err := sd.GetCellRichText(sheetName, srcCell)
					if err == nil && len(richText) > 0 {
						f.SetCellRichText(sheetName, dstCell, richText)
					} else {
						f.SetCellValue(sheetName, dstCell, rows[c2Row][0])
					}
					copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
				}

				log.Printf("%s: Запись данных для latex завершена", sheetName)
			}
		}

		// Устанавливаем ширину колонок
		f.SetColWidth(sheetName, "A", "A", 23.14)
		f.SetColWidth(sheetName, "F", "F", 15)

		f.SetActiveSheet(index)
	}

	// Переименовываем первый лист
	err = f.SetSheetName("Sheet1", "пробоподготовка")
	if err != nil {
		log.Println("Ошибка переименовывания листа на \"пробоподготовка\"", err)
	}

	// Добавляем выводы
	_, err = f.NewSheet("выводы")
	if err != nil {
		log.Print("Ошибка выводов: ", err)
	}

	// Сохраняем файл
	err = f.SaveAs(resultFileName)
	if err != nil {
		log.Print("Ошибка сохранения: ", err)
	}
}

// Копирование стиля между файлами
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

	dst.SetCellStyle(dstSheet, dstCell, dstCell, newStyleID)
}
