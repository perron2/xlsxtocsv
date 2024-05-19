ifdef ComSpec
	EXE := .exe
endif

build:
	go build -ldflags="-w -s" -o bin/xlsxtocsv$(EXE) perron2.ch/xlsxtocsv

format:
	go fmt ./...
