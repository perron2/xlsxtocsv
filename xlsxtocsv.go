package main

import (
	"encoding/csv"
	"fmt"
	"github.com/thedatashed/xlsxreader"
	"io"
	"os"
	"path/filepath"
	"perron2.ch/xlsxtocsv/config"
	"strings"
)

const (
	programVersion = "1.0"
)

func main() {
	cfg := config.Read(programVersion)
	for _, inputFile := range cfg.InputFiles {
		convertFile(inputFile, cfg)
	}
}

func convertFile(inputFile string, cfg config.Config) {
	xl, err := xlsxreader.OpenFile(inputFile)
	if err != nil {
		errorExit("Cannot open Excel file \"%s\": %s", inputFile, err)
	}
	defer xl.Close()

	outputPath := cfg.OutputFile
	if cfg.OutputDir != "" {
		fileName := filepath.Base(inputFile)
		fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".csv"
		outputPath = filepath.Join(cfg.OutputDir, fileName)
	}

	outputFile := os.Stdout
	if outputPath != "" {
		outputFile, err = os.Create(outputPath)
		if err != nil {
			errorExit("Cannot create output file \"%s\": %s", outputPath, err)
		}
		defer outputFile.Close()
	}

	encodedWriter := io.Writer(outputFile)
	if cfg.Charset == config.Ansi {
		encodedWriter = ansiWriter{w: outputFile, fallback: '?'}
	}

	csvWriter := csv.NewWriter(encodedWriter)
	csvWriter.Comma = cfg.Separator
	defer csvWriter.Flush()

	first := true
	columnCount := -1
	for row := range xl.ReadRows(xl.Sheets[0]) {
		var record []string
		for _, cell := range row.Cells {
			index := cell.ColumnIndex()
			for len(record) < index {
				record = append(record, "")
			}
			record = append(record, cell.Value)
		}
		if first {
			first = false
			if !cfg.Headers {
				continue
			}
			columnCount = len(record)
			fileMappings := cfg.FileMappings[filepath.Base(inputFile)]
			for i, header := range record {
				if mapped, ok := fileMappings[header]; ok {
					record[i] = mapped
				} else if mapped, ok := cfg.GlobalMappings[header]; ok {
					record[i] = mapped
				}
			}
		}
		if cfg.Headers {
			for len(record) < columnCount {
				record = append(record, "")
			}
		}
		if err := csvWriter.Write(record); err != nil {
			out := "standard output"
			if outputPath != "" {
				out = "\"" + outputPath + "\""
			}
			errorExit("Error while writing to %s: %s", out, err)
		}
	}
}

func errorExit(format string, a ...any) {
	fmt.Printf(format, a...)
	os.Exit(1)
}
