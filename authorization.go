package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"log"
	"os"
)

var (
	key = flag.String("key", "", "key string")
)

var keys = []string{
	"fba9c56168b64dfd2803a681ae1e83c1",
	"bda657d55d3315684a2b3ce453d02891",
}

func authorization() {
	log.Println("check the permission.")

	nameTranslated := md5Translate(parseKey())
	crackDebug("translated key:%v\n", nameTranslated)
	for _, key := range keys {
		if key == nameTranslated {
			log.Println("permission check successfully.")
			return
		}
	}
	log.Println("permission check failed.")
	os.Exit(0)
}
func parseKey() string {
	flag.Parse()
	crackDebug("input key:%s\n", *key)
	return *key
}
func md5Translate(name string) string {
	h := md5.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}
