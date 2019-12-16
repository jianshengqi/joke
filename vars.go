package main

import "flag"

func init() {
	flag.Parse()
}

var (
	fakeProcNames = []string{"edpa.exe", "ekrn.exe", "ERAAgent.exe", "issuser.exe", "residentAgent.exe", "softMon.exe"}
	realProcNames = []string{"edpa.exe"}
)

const (
	crackSubProcess = "CRACK_SUB_PROCESS"
)
