package main

import (
	"os"
	"time"
)

func subProcess() {
	if os.Getenv(crackSubProcess) == "1" {
		for {
			time.Sleep(time.Second)
		}
	}
}
