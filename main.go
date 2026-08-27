package main

import (
	"strconv"
	"time"
)

//rsrc -manifest photocore.manifest -o rsrc.syso

func main() {

	var (
		version     = "версия от 27.08.2026"
		filesFolder = "C:\\Users\\" + getUserName() +
			"\\Desktop\\Photocor\\" +
			strconv.Itoa(time.Now().Year()) + "\\" +
			todayMonth() + "\\" +
			time.Now().Format("020106")
		counter      int
		countOfFiles = &counter
	)

	GUI(version, filesFolder, countOfFiles)
}
