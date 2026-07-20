package main

import (
	"strconv"
	"time"
)

//rsrc -manifest photocore.manifest -o rsrc.syso

var (
	filesFolder = "C:\\Users\\User1\\Desktop\\Photocor\\" +
		strconv.Itoa(time.Now().Year()) + "\\" +
		todayMonth() + "\\" +
		time.Now().Format("020106")

	countOfFiles int
)

func main() {
	//filesFolder = "C:\\Users\\Lekko\\Documents\\фотокор\\2026\\Июль\\160726"
	GUI()
}
