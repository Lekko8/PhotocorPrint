package main

import "strings"

func buildNames(files []string) string {
	if len(files) == 0 {
		return ""
	}
	var filenames strings.Builder
	for _, file := range files {
		if filenames.Len() != 0 {
			filenames.WriteString("\r\n")
		}
		filenames.WriteString(file)
	}
	return filenames.String()
}

func createFile(resultFileName string) {

}
