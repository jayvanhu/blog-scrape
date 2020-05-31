package util

import (
	"os"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func WriteToFile(filename, text string) {
	file, err := os.Create(filename)
	HandleErr(err)
	_, err = file.WriteString(text)
	HandleErr(err)
	file.Close()
}

func AppendToFile(filename, text string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	HandleErr(err)
	_, err = file.WriteString(text)
	HandleErr(err)
	file.Close()
}
