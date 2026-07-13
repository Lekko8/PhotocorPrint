package main

import (
	"log"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	GUI()
	log.Println("Приложение запущено")
}

func GUI() {

	var (
		inputData      *walk.LineEdit
		inputOrderName *walk.LineEdit
		statusFiles    *walk.TextEdit
		statusLabel    *walk.Label
		mw             *walk.MainWindow
	)

	targetDate := time.Now().Format("02.01.2006")

	var capturedDate = targetDate
	//appIcon, _ := walk.NewIconFromResourceId(2)

	mainWindow := MainWindow{
		AssignTo: &mw,
		Title:    "PhotocoreData",
		//Icon:     appIcon,
		Size:   Size{Width: 350, Height: 400},
		Layout: VBox{},
		Children: []Widget{
			Label{Text: "Дата искомых файлов:"},
			LineEdit{AssignTo: &inputData},

			PushButton{
				Text: "Повторить поиск",
				OnClicked: func() {
					capturedDate = inputData.Text()
					statusLabel.SetText("Прочитаны файлы за " + capturedDate)
				},
			},

			Label{Text: "Найденные файлы:"},
			TextEdit{AssignTo: &statusFiles, ReadOnly: true, VScroll: true},

			Label{Text: "Введите номер заказа:"},
			LineEdit{AssignTo: &inputOrderName},

			PushButton{
				Text: "Создать .xlsx",
				OnClicked: func() {
					//orderName := inputOrderName.Text()
					//files := dataCreate(capturedFiles)
					//results := calculateFiles(files)
					//createrXlsx(files, results, orderName, capturedDate)
					statusLabel.SetText("Файл .xlsx успешно создан")
				},
			},

			Label{AssignTo: &statusLabel},
		},
	}

	if err := mainWindow.Create(); err != nil {
		log.Fatal(err) // Увидите реальную ошибку
	}
	mw.Run()
}
