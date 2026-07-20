package main

import (
	_ "image/jpeg"
	"log"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

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

	picsCellsF := []string{"G13", "N13", "N24"} // Целевые ячейки для графиков

	for idx := range groupsList {
		sheetName := groupsList[idx]

		index, err := f.NewSheet(sheetName)
		if err != nil {
			log.Printf("Ошибка создания листа %s: %v", sheetName, err)
			continue
		}

		m, err := sd.GetMergeCells(sheetName, true)
		if err != nil {
			log.Print(err.Error())
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

		picsCellSD := []string{p1, p2, p3} // Копируем графики

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
					return "Ошибка вставки картинки: " + err.Error()
				}
			}
		}

		rows, err := sd.GetRows(sheetName) // Читаем всю страницу
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

		currentRow := 10 // Отступаем шапку

		for rowNum := 1; rowNum < startReadingRow; rowNum++ {

			srcCell, err := excelize.CoordinatesToCellName(1, rowNum)
			if err != nil {
				log.Print(err.Error())
				return err.Error()
			}

			val, err := sd.GetCellValue(sheetName, srcCell)
			if err != nil {
				continue
			}

			if val == "" {
				continue
			}

			richText, err := sd.GetCellRichText(sheetName, srcCell)
			if err != nil {
				log.Print(err.Error())
				return err.Error()
			}

			dstCell, err := excelize.CoordinatesToCellName(1, currentRow)
			if err != nil {
				log.Print(err.Error())
				return err.Error()
			}

			if len(richText) > 0 {
				err = f.SetCellRichText(sheetName, dstCell, richText)
				if err != nil {
					log.Print(err.Error())
					return err.Error()
				}
			} else {
				err = f.SetCellValue(sheetName, dstCell, val)
				if err != nil {
					log.Print(err.Error())
					return err.Error()
				}
			}

			copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
			currentRow++
		}

		rowsAfterCopy, err := f.GetRows(sheetName)
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
		for i, row := range rowsAfterCopy {
			if len(row) > 0 && strings.Contains(row[0], "File(s):") {
				err = f.InsertRows(sheetName, i+2, 1) // Добавляем пустую строку после "File(s): ..."
				if err != nil {
					log.Print(err.Error())
					return err.Error()
				}
				break
			}
		}

		headerEndRow := currentRow - 1 // Запоминаем строку, где закончились маркеры

		if headerEndRow < 19 {
			headerEndRow = currentRow - 1
		}

		var headerRows []int                                          // Начинаем копировать остальные данные
		for rowNum := startReadingRow; rowNum < len(rows); rowNum++ { // Находим блок заголовков
			if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "Distribution analysis") {

				for j := rowNum - 1; j <= rowNum+5 && j < len(rows); j++ {
					if j >= startReadingRow {
						headerRows = append(headerRows, j)
					}
				}
				break

			}
		}

		type Block struct {
			markerRow  int
			tableStart int
			tableEnd   int
			c2Row      int
		}
		var blocks []Block // Находим все блоки с данными

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

		dataStartRow := headerEndRow + 2 // Вставляем данные после имён файлов
		currentRow = dataStartRow

		if len(blocks) > 0 { // Записываем заголовки через RichText

			for _, headerRow := range headerRows {
				if headerRow < len(rows) {
					for colIdx := 0; colIdx < 5 && colIdx < len(rows[headerRow]); colIdx++ {
						if rows[headerRow][colIdx] != "" {
							srcCell, err := excelize.CoordinatesToCellName(colIdx+1, headerRow+1)
							if err != nil {
								log.Print(err.Error())
								return err.Error()
							}
							dstCell, err := excelize.CoordinatesToCellName(colIdx+1, currentRow)
							if err != nil {
								log.Print(err.Error())
								return err.Error()
							}

							richText, err := sd.GetCellRichText(sheetName, srcCell)
							if err == nil && len(richText) > 0 {
								err = f.SetCellRichText(sheetName, dstCell, richText)
								if err != nil {
									log.Print(err.Error())
									return err.Error()
								}
							} else {
								err = f.SetCellValue(sheetName, dstCell, rows[headerRow][colIdx])
								if err != nil {
									log.Print(err.Error())
									return err.Error()
								}
							}
							copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
						}
					}
					currentRow++
				}
			}

			currentRow++ // Пустая строка после заголовков

			for _, block := range blocks {

				if block.markerRow < len(rows) {
					srcCell, err := excelize.CoordinatesToCellName(1, block.markerRow+1)
					if err != nil {
						log.Print(err.Error())
						return err.Error()
					}
					dstCell, err := excelize.CoordinatesToCellName(1, currentRow)
					if err != nil {
						log.Print(err.Error())
						return err.Error()
					}

					richText, err := sd.GetCellRichText(sheetName, srcCell)
					if err == nil && len(richText) > 0 {
						err = f.SetCellRichText(sheetName, dstCell, richText)
						if err != nil {
							log.Print(err.Error())
							return err.Error()
						}
					} else {
						err = f.SetCellValue(sheetName, dstCell, rows[block.markerRow][0])
						if err != nil {
							log.Print(err.Error())
							return err.Error()
						}
					}
					copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
					currentRow++
				}

				currentRow++ // Пустая строка

				for rowIdx := block.tableStart; rowIdx <= block.tableEnd && rowIdx < len(rows); rowIdx++ {
					if len(rows[rowIdx]) > 0 {
						for colIdx := 0; colIdx < len(rows[rowIdx]) && colIdx < 5; colIdx++ {
							if rows[rowIdx][colIdx] != "" {
								srcCell, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
								if err != nil {
									log.Print(err.Error())
									return err.Error()
								}
								dstCell, err := excelize.CoordinatesToCellName(colIdx+1, currentRow)
								if err != nil {
									log.Print(err.Error())
									return err.Error()
								}

								richText, err := sd.GetCellRichText(sheetName, srcCell)
								if err == nil && len(richText) > 0 && strings.Contains(rows[rowIdx][colIdx], "c") {
									err = f.SetCellRichText(sheetName, dstCell, richText)
									if err != nil {
										log.Print(err.Error())
										return err.Error()
									}
								} else {
									err = f.SetCellValue(sheetName, dstCell, rows[rowIdx][colIdx])
									if err != nil {
										log.Print(err.Error())
										return err.Error()
									}
								}

								copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
							}
						}
						currentRow++
					}
				}

				if block.c2Row > 0 && block.c2Row < len(rows) && len(rows[block.c2Row]) > 0 {
					c2Row := currentRow - 1
					srcCell, err := excelize.CoordinatesToCellName(1, block.c2Row+1)
					if err != nil {
						log.Print(err.Error())
						return err.Error()
					}
					dstCell, err := excelize.CoordinatesToCellName(6, c2Row)
					if err != nil {
						log.Print(err.Error())
						return err.Error()
					}

					richText, err := sd.GetCellRichText(sheetName, srcCell)
					if err == nil && len(richText) > 0 {
						err = f.SetCellRichText(sheetName, dstCell, richText)
						if err != nil {
							log.Print(err.Error())
							return err.Error()
						}
					} else {
						err = f.SetCellValue(sheetName, dstCell, rows[block.c2Row][0])
						if err != nil {
							log.Print(err.Error())
							return err.Error()
						}
					}
					copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
				}

				currentRow++ // Пустая строка после блока
			}

		} else { // Если только 1 файл

			var fileMarker string // Находим маркер File:
			var fileMarkerRow int
			for rowNum := startReadingRow; rowNum < len(rows); rowNum++ {
				if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "File: ") {
					fileMarker = rows[rowNum][0]
					fileMarkerRow = rowNum
					break
				}
			}

			if fileMarker == "" { // Если маркер не найден, ищем в строках до startReadingRow
				for rowNum := 0; rowNum < startReadingRow; rowNum++ {
					if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "File: ") {
						fileMarker = rows[rowNum][0]
						fileMarkerRow = rowNum
						break
					}
				}
			}

			var dataRows []int // Находим данные
			for rowNum := startReadingRow; rowNum < len(rows); rowNum++ {
				if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "Distribution analysis") {
					for j := rowNum - 1; j <= rowNum+5 && j < len(rows); j++ {
						if j >= startReadingRow {
							dataRows = append(dataRows, j)
						}
					}
					break
				}
			}

			var tableStart, tableEnd, c2Row int // Ищем таблицу
			for rowNum := startReadingRow; rowNum < len(rows); rowNum++ {
				if len(rows[rowNum]) > 0 && strings.Contains(rows[rowNum][0], "Peak Num") {
					tableStart = rowNum
					for j := tableStart + 1; j < len(rows); j++ {
						if len(rows[j]) == 0 {
							continue
						}
						if strings.HasPrefix(rows[j][0], "c") {
							c2Row = j
							tableEnd = j - 1
							break
						}
					}
					break
				}
			}

			if tableStart > 0 && fileMarker != "" {

				currentRow++ // Добавляем пустую строку

				srcCell, err := excelize.CoordinatesToCellName(1, fileMarkerRow+1)
				if err != nil {
					log.Print(err.Error())
					return err.Error()
				}
				dstCell, err := excelize.CoordinatesToCellName(1, currentRow)
				if err != nil {
					log.Print(err.Error())
					return err.Error()
				}

				richText, err := sd.GetCellRichText(sheetName, srcCell)
				if err == nil && len(richText) > 0 {
					err = f.SetCellRichText(sheetName, dstCell, richText)
					if err != nil {
						log.Print(err.Error())
						return err.Error()
					}
				} else {
					err = f.SetCellValue(sheetName, dstCell, fileMarker)
					if err != nil {
						log.Print(err.Error())
						return err.Error()
					}
				}
				copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
				currentRow++

				currentRow++ // Добавляем пустую строку

				for _, rowIdx := range dataRows { // Записываем заголовки
					if rowIdx < len(rows) {
						for colIdx := 0; colIdx < 5 && colIdx < len(rows[rowIdx]); colIdx++ {
							if rows[rowIdx][colIdx] != "" {
								srcCell, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
								if err != nil {
									log.Print(err.Error())
									return err.Error()
								}
								dstCell, err := excelize.CoordinatesToCellName(colIdx+1, currentRow)
								if err != nil {
									log.Print(err.Error())
									return err.Error()
								}

								richText, err := sd.GetCellRichText(sheetName, srcCell)
								if err == nil && len(richText) > 0 {
									err = f.SetCellRichText(sheetName, dstCell, richText)
									if err != nil {
										log.Print(err.Error())
										return err.Error()
									}
								} else {
									err = f.SetCellValue(sheetName, dstCell, rows[rowIdx][colIdx])
									if err != nil {
										log.Print(err.Error())
										return err.Error()
									}
								}
								copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
							}
						}
						currentRow++
					}
				}

				currentRow++ // Добавляем пустую строку

				for rowIdx := tableStart; rowIdx <= tableEnd && rowIdx < len(rows); rowIdx++ {
					if len(rows[rowIdx]) > 0 {
						for colIdx := 0; colIdx < len(rows[rowIdx]) && colIdx < 5; colIdx++ {
							if rows[rowIdx][colIdx] != "" {
								srcCell, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
								if err != nil {
									log.Print(err.Error())
									return err.Error()
								}
								dstCell, err := excelize.CoordinatesToCellName(colIdx+1, currentRow)
								if err != nil {
									log.Print(err.Error())
									return err.Error()
								}

								richText, err := sd.GetCellRichText(sheetName, srcCell)
								if err == nil && len(richText) > 0 &&
									strings.Contains(rows[rowIdx][colIdx], "c") {

									err = f.SetCellRichText(sheetName, dstCell, richText)
									if err != nil {
										log.Print(err.Error())
										return err.Error()
									}

								} else {
									err = f.SetCellValue(sheetName, dstCell, rows[rowIdx][colIdx])
									if err != nil {
										log.Print(err.Error())
										return err.Error()
									}
								}
								copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
							}
						}
						currentRow++
					}
				}

				if c2Row > 0 && c2Row < len(rows) && len(rows[c2Row]) > 0 {
					c2RowDst := currentRow - 1
					srcCell, err := excelize.CoordinatesToCellName(1, c2Row+1)
					if err != nil {
						log.Print(err.Error())
						return err.Error()
					}
					dstCell, err := excelize.CoordinatesToCellName(6, c2RowDst)
					if err != nil {
						log.Print(err.Error())
						return err.Error()
					}

					richText, err := sd.GetCellRichText(sheetName, srcCell)
					if err == nil && len(richText) > 0 {
						err = f.SetCellRichText(sheetName, dstCell, richText)
						if err != nil {
							log.Print(err.Error())
							return err.Error()
						}
					} else {
						err = f.SetCellValue(sheetName, dstCell, rows[c2Row][0])
						if err != nil {
							log.Print(err.Error())
							return err.Error()
						}
					}
					copyCellStyle(sd, f, sheetName, sheetName, srcCell, dstCell)
				}
			}
		}

		err = f.SetCellValue(sheetName, "D13", "Mean intensity")
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
		err = f.SetCellValue(sheetName, "E13", "Rate of correct unit")
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}

		for i := 0; true; i++ {

			contCell, err := f.GetCellValue(sheetName, "A"+strconv.Itoa(i+14))
			if err != nil {
				log.Print(err.Error())
				return err.Error()
			}

			if contCell == "" {
				break
			}

			contCell = strings.TrimPrefix(contCell, "¨ ")

			err = f.SetCellValue(sheetName, "D"+strconv.Itoa(i+14), reportMatrix[contCell].MeanIntensity)
			if err != nil {
				log.Print(err.Error())
				return err.Error()
			}
			err = f.SetCellValue(sheetName, "E"+strconv.Itoa(i+14), reportMatrix[contCell].RateOfCorrectUnit)
			if err != nil {
				log.Print(err.Error())
				return err.Error()
			}

		}

		err = f.SetColWidth(sheetName, "A", "A", 23.14)
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
		err = f.SetColWidth(sheetName, "D", "F", 15)
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}

		f.SetActiveSheet(index)

	}

	err = f.SetSheetName("Sheet1", "пробоподготовка")
	if err != nil {
		log.Println("Ошибка переименовывания листа на \"пробоподготовка\"", err)
		return "Ошибка переименовывания листа на \"пробоподготовка\"" + err.Error()
	}

	_, err = f.NewSheet("выводы")
	if err != nil {
		log.Print("Ошибка выводов: ", err)
		return "Ошибка выводов: " + err.Error()
	}

	err = f.SaveAs(resultFileName)
	if err != nil {
		log.Print("Ошибка сохранения: ", err)
		return "Ошибка сохранения: " + err.Error()
	}

	return "Файл " + resultFileName + " успешно создан"
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

	_ = dst.SetCellStyle(dstSheet, dstCell, dstCell, newStyleID)
}
