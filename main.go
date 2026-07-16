package main

import (
	"strconv"
	"time"
)

//rsrc -manifest photocore.manifest -o rsrc.syso

//"C:\\Users\\user\\Desktop\\Photocor\\"
//C:\Users\Lekko\Documents\фотокор\2025\Июнь\240625

var filesFolder = "C:\\Users\\Lekko\\Documents\\фотокор\\" +
	strconv.Itoa(time.Now().Year()) + "\\" +
	todayMonth() + "\\" +
	time.Now().Format("020106")
var countOfFiles int

func main() {
	GUI()
}
