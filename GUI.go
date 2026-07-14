package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func GUI() {

	var (
		inputData      *walk.LineEdit
		inputOrderName *walk.LineEdit
		statusFiles    *walk.TextEdit
		statusLabel    *walk.Label
		mw             *walk.MainWindow
	)

	//appIcon, _ := walk.NewIconFromResourceId(2)

	mainWindow := MainWindow{
		AssignTo: &mw,
		Title:    "DLS Фотокор",
		//Icon:     appIcon,
		Size:   Size{Width: 350, Height: 400},
		Layout: VBox{},
		Children: []Widget{
			Label{Text: "Путь к папке:"},
			LineEdit{AssignTo: &inputData, Text: filesFolder},

			PushButton{
				Text: "Повторить поиск",
				OnClicked: func() {
					err := statusLabel.SetText(search(inputData.Text()))
					if err != nil {
						log.Panic(err.Error())
						return
					}
				},
			},

			Label{Text: "Найденные файлы:"},
			TextEdit{AssignTo: &statusFiles, ReadOnly: true, VScroll: true},

			Label{Text: "Инициалы исполнителя:"},
			LineEdit{AssignTo: &inputOrderName},

			PushButton{
				Text: "Создать .xlsx",
				OnClicked: func() {
					err := statusLabel.SetText(createXlsx(inputOrderName.Text()))
					if err != nil {
						log.Panic(err.Error())
					}
				},
			},

			Label{AssignTo: &statusLabel, Text: "Лёша лох"},
		},
	}

	if err := mainWindow.Create(); err != nil {
		log.Fatal(err)
	}

	statusLabel.SetText(search(filesFolder)) //запуск поиска на старте

	mw.Run()
}
