package main

import (
	"fmt"
	"log"
	"os"
)

func crackInfo(format string, v ...interface{}) {
	log.Printf(format, v...)
}
func crackDebug(format string, v ...interface{}) {
	if os.Getenv("CRACK_DEBUG") != "1" {
		return
	}
	log.Printf(format, v...)
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
