package main

import (
	"strconv"
	"time"
)

//rsrc -manifest photocore.manifest -o rsrc.syso

//"C:\\Users\\user\\Desktop\\Photocor\\"

var filesFolder = "C:\\Users\\Lekko\\Documents\\фотокор\\" +
	strconv.Itoa(time.Now().Year()) + "\\" +
	todayMonth() + "\\" +
	time.Now().Format("020106")
var countOfFiles int

func main() {
	filesFolder = "C:\\Users\\Lekko\\Documents\\фотокор\\2026\\Июль\\160726"
	GUI()
}
