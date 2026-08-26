package main

import (
	"log"
	"os/user"
	"strconv"
	"strings"
	"time"
)

//rsrc -manifest photocore.manifest -o rsrc.syso

func getUserName() string {
	corUse, err := user.Current()
	if err != nil {
		log.Print(err.Error())
		return ""
	}
	return strings.Split(corUse.Username, "\\")[1]
}

var filesFolder = "C:\\Users\\" + getUserName() +
	"\\Desktop\\Photocor\\" +
	strconv.Itoa(time.Now().Year()) + "\\" +
	todayMonth() + "\\" +
	time.Now().Format("020106")
var countOfFiles int

func main() {
	GUI()
}
