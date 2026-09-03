package main

import (
	"log"
	"strconv"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func GUI(version, filesFolder string, countOfFiles *int) {

	var (
		inputFolder  *walk.LineEdit
		inputOrder   *walk.LineEdit
		inputName    *walk.LineEdit
		inputMin     *walk.LineEdit
		inputMax     *walk.LineEdit
		inputSec     *walk.LineEdit
		statusFiles  *walk.TextEdit
		statusSearch *walk.Label
		statusProg   *walk.TextEdit
		mw           *walk.MainWindow
	)

	appIcon, _ := walk.NewIconFromResourceId(2)

	mainWindow := MainWindow{
		AssignTo: &mw,
		Title:    "DLS Фотокор | " + version,
		Icon:     appIcon,
		Size:     Size{Width: 400, Height: 400},
		MinSize:  Size{Width: 400, Height: 400},
		Layout:   VBox{},
		Children: []Widget{

			HSplitter{
				Children: []Widget{
					Label{Text: "Путь к папке:"},
					LineEdit{
						AssignTo: &inputFolder,
						Text:     filesFolder,
						MaxSize:  Size{Width: 300, Height: 25},
					},
				},
			},

			HSplitter{
				Children: []Widget{
					PushButton{
						MinSize: Size{Width: 30, Height: 25},
						MaxSize: Size{Width: 120, Height: 25},
						Text:    "Повторить поиск",
						OnClicked: func() {
							*countOfFiles = 0
							filesFolder = inputFolder.Text()

							err := statusProg.SetText("Прочитаны файлы из:\r\n" + filesFolder)
							if err != nil {
								log.Panic(err.Error())
							}

							err = statusFiles.SetText(readFileList(filesFolder, *countOfFiles))
							if err != nil {
								log.Panic(err.Error())
							}

							err = statusSearch.SetText("Найденные файлы: " + strconv.Itoa(*countOfFiles))
							if err != nil {
								log.Panic(err.Error())
							}
						},
					},
					Label{Text: "Заказ:"},
					LineEdit{
						AssignTo: &inputOrder,
						MinSize:  Size{Width: 80, Height: 25},
						MaxSize:  Size{Width: 80, Height: 25},
					},
					Label{Text: "Инициалы:"},
					LineEdit{
						AssignTo: &inputName,
						MinSize:  Size{Width: 50, Height: 25},
						MaxSize:  Size{Width: 50, Height: 25},
					},
				},
			},

			HSplitter{
				Children: []Widget{
					//TextLabel{Text: "Критерии оценки пиков:"},
					TextLabel{Text: "Минимум"},
					LineEdit{
						AssignTo: &inputMin,
						Text:     "2.9",
						//MinSize:  Size{Width: 80, Height: 25},
						//MaxSize:  Size{Width: 80, Height: 25},
					},
					TextLabel{Text: "Максимум"},
					LineEdit{
						AssignTo: &inputMax,
						Text:     "9.0",
						//MinSize:  Size{Width: 80, Height: 25},
						//MaxSize:  Size{Width: 80, Height: 25},
					},
					TextLabel{Text: "Второй пик"},
					LineEdit{
						AssignTo: &inputSec,
						Text:     "30",
						//MinSize:  Size{Width: 80, Height: 25},
						//MaxSize:  Size{Width: 80, Height: 25},
					},
				},
			},

			HSplitter{
				Children: []Widget{
					PushButton{
						MaxSize: Size{Width: 100, Height: 25},
						Text:    "Создать .xlsx",
						OnClicked: func() {
							err := statusProg.SetText("Идёт сборка файла .xlsx")
							if err != nil {
								log.Panic(err.Error())
							}
							filesFolder = inputFolder.Text()
							order := inputOrder.Text()
							name := inputName.Text()
							err = statusProg.SetText(createFile(filesFolder, order, name))
							if err != nil {
								log.Panic(err.Error())
							}
						},
					},
					PushButton{
						MaxSize: Size{Width: 100, Height: 25},
						Text:    "Обсчитать файлик",
						OnClicked: func() {
							err := statusProg.SetText("Обсчитываю... хихи")
							if err != nil {
								log.Panic(err.Error())
							}
							order := inputOrder.Text()
							name := inputName.Text()
							KMin := inputMin.Text()
							KMax := inputMax.Text()
							KSec := inputSec.Text()
							err = statusProg.SetText(calculate(createFileName(order, name), KMin, KMax, KSec))
							if err != nil {
								log.Panic(err.Error())
							}
						},
					},
				},
			},

			Label{AssignTo: &statusSearch, Text: "Найденные файлы:"},
			TextEdit{AssignTo: &statusFiles, ReadOnly: true, VScroll: true},
			TextEdit{
				AssignTo: &statusProg,
				Text:     "Лёша лох",
				MaxSize:  Size{Width: 300, Height: 40},
				ReadOnly: true,
				VScroll:  true,
			},
		},
	}

	if err := mainWindow.Create(); err != nil {
		log.Fatal(err)
	}

	err := statusProg.SetText("Прочитаны файлы из:\r\n" + filesFolder) //запуск поиска на старте
	if err != nil {
		log.Panic(err.Error())
	}
	err = statusFiles.SetText(readFileList(filesFolder, *countOfFiles))
	if err != nil {
		log.Panic(err.Error())
	}
	err = statusSearch.SetText("Найденные файлы: " + strconv.Itoa(*countOfFiles))
	if err != nil {
		log.Panic(err.Error())
	}

	mw.Run()
}
