BUILDFLAGS := -ldflags="-s -w" -trimpath

ifdef ComSpec
	EXE := .exe
endif

build:
	go build -o bin/xlsxtocsv$(EXE) perron2.ch/xlsxtocsv

windows:
	go build $(BUILDFLAGS) -o bin/windows/xlsxtocsv.exe perron2.ch/xlsxtocsv
	upx bin/windows/xlsxtocsv.exe

darwin:
	GOOS=darwin GOARCH=amd64 go build $(BUILDFLAGS) -o bin/darwin/xlsxtocsv perron2.ch/xlsxtocsv

linux:
	GOOS=linux GOARCH=amd64 go build $(BUILDFLAGS) -o bin/linux/xlsxtocsv perron2.ch/xlsxtocsv
	upx bin/linux/xlsxtocsv

format:
	go fmt ./...
