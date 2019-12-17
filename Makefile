.PHONY: build syso pack upx

bin=joke.exe

upx:build
	upx $(bin)

build:clean syso pack
	go build -o $(bin) -ldflags "-s -w"

syso:
	windres -o joke.syso joke.rc

pack:
	go-bindata -o=edpa.go edpa.exe

run:
	$(bin) -key taozhang8

clean:
	-del $(bin)