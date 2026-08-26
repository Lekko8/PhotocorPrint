package main

import (
	"log"
	"strconv"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func GUI() {

	var (
		inputFolder  *walk.LineEdit
		inputName    *walk.LineEdit
		statusFiles  *walk.TextEdit
		statusSearch *walk.Label
		statusProg   *walk.TextEdit
		mw           *walk.MainWindow
	)

	appIcon, _ := walk.NewIconFromResourceId(2)

	mainWindow := MainWindow{
		AssignTo: &mw,
		Title:    "DLS Фотокор | версия от 26.08.2026",
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
							countOfFiles = 0
							filesFolder = inputFolder.Text()

							err := statusProg.SetText("Прочитаны файлы из:\r\n" + filesFolder)
							if err != nil {
								log.Panic(err.Error())
							}

							err = statusFiles.SetText(readFileList(filesFolder))
							if err != nil {
								log.Panic(err.Error())
							}

							err = statusSearch.SetText("Найденные файлы: " + strconv.Itoa(countOfFiles))
							if err != nil {
								log.Panic(err.Error())
							}
						},
					},
					Label{Text: "Инициалы (в конец имени файла):"},
					LineEdit{
						AssignTo: &inputName,
						MinSize:  Size{Width: 50, Height: 25},
						MaxSize:  Size{Width: 50, Height: 25},
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
							err = statusProg.SetText(createFile())
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
							name := inputName.Text()
							err = statusProg.SetText(calculate(createFileName(name)))
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
	err = statusFiles.SetText(readFileList(filesFolder))
	if err != nil {
		log.Panic(err.Error())
	}
	err = statusSearch.SetText("Найденные файлы: " + strconv.Itoa(countOfFiles))
	if err != nil {
		log.Panic(err.Error())
	}

	mw.Run()
}
