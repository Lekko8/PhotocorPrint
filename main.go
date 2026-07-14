package main

import (
	"strconv"
	"time"
)

//rsrc -manifest photocore.manifest -o rsrc.syso

var filesFolder = "C:\\Users\\user\\Desktop\\Photocor\\" +
	strconv.Itoa(time.Now().Year()) + "\\" +
	todayMonth() + "\\" +
	time.Now().Format("020106")

func main() {
	GUI()
}
