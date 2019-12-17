package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func crackInfo(format string, v ...interface{}) {
	log.Printf(format, v...)
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
func pause() {
	time.Sleep(time.Second)
}
