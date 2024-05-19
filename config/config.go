package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

type Charset int8

const (
	Utf8 Charset = iota
	Ansi
)

type Config struct {
	Headers        bool
	Separator      rune
	Charset        Charset
	GlobalMappings map[string]string
	FileMappings   map[string]map[string]string
	InputFiles     []string
	OutputFile     string
	OutputDir      string
}

func Read(programVersion string) Config {
	showVersion := flag.Bool("version", false, "Show program version number")
	noHeaders := flag.Bool("noheaders", false, "Do not generate a header line")
	separator := flag.String("separator", ",", "Field separator")
	charset := flag.String("charset", "utf8", "CSV file character set encoding")
	mapFile := flag.String("mapfile", "", "Map file with mapping specifications")
	outputFile := flag.String("out", "", "Output CSV file (default \"stdout\")")
	outputDir := flag.String("outdir", "", "Output directory for CSV files")

	var mappings stringArrayFlag
	flag.Var(&mappings, "map", "Map an Excel column name to a CSV column name (from=to)")

	programName := filepath.Base(os.Args[0])
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options] <input file> [<input file>...]\n\nOptions:\n", programName)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("%s v%s\n", programName, programVersion)
		os.Exit(0)
	}

	config := Config{
		Headers:        !*noHeaders,
		Separator:      readSeparator(*separator),
		Charset:        readCharset(*charset),
		GlobalMappings: make(map[string]string),
		FileMappings:   make(map[string]map[string]string),
		OutputFile:     *outputFile,
		OutputDir:      *outputDir,
	}

	if config.OutputFile != "" && config.OutputDir != "" {
		errorExit("Specifying both an output file and an output directory is not allowed")
	} else if config.OutputDir != "" {
		if stat, err := os.Stat(config.OutputDir); err == nil && !stat.IsDir() {
			errorExit("\"%s\" is not a directory", config.OutputDir)
		} else if os.IsNotExist(err) {
			errorExit("Directory \"%s\" does not exist", config.OutputDir)
		} else if err != nil && stat.IsDir() {
			errorExit("Cannot check directory  \"%s\": %ys", config.OutputDir, err)
		}
	}

	config.InputFiles = readInputFiles(flag.Args())
	if len(config.InputFiles) == 0 {
		errorExit("No input file specified")
	} else if len(config.InputFiles) > 1 && config.OutputDir == "" {
		errorExit("Multiple input files (%d) cannot be written to a single output file or to standard output", len(config.InputFiles))
	}

	if *mapFile != "" {
		readMapFile(*mapFile, config.GlobalMappings, config.FileMappings)
	}

	readMappings(mappings, config.GlobalMappings)

	return config
}

func readInputFiles(args []string) []string {
	var inputFiles []string
	for _, arg := range args {
		if stat, err := os.Stat(arg); err == nil && !stat.IsDir() {
			inputFiles = append(inputFiles, arg)
			continue
		}
		matches, err := filepath.Glob(arg)
		if err != nil {
			errorExit("Cannot analyze input \"%s\": %s", arg, err)
		} else if len(matches) == 0 {
			if strings.ContainsAny(arg, "*?[]") {
				errorExit("\"%s\" does not match anything", arg)
			} else {
				errorExit("\"%s\" does not exist", arg)
			}
		}
		inputFiles = append(inputFiles, matches...)
	}
	slices.Sort(inputFiles)
	return slices.Compact(inputFiles)
}

func readSeparator(separator string) rune {
	switch separator {
	case ",", ";":
	default:
		errorExit("Invalid separator; only \",\" and \";\" are supported")
	}
	return []rune(separator)[0]
}

func readCharset(charset string) Charset {
	switch strings.ToLower(charset) {
	case "", "utf8", "utf-8":
		return Utf8
	case "ansi":
		return Ansi
	default:
		errorExit("Invalid character set %s; supported character sets are \"utf8\" and \"ansi\"")
		return Utf8
	}
}

func readMappings(mappings []string, m map[string]string) {
	for _, mapping := range mappings {
		parts := strings.Split(mapping, "=")
		if len(parts) != 2 {
			errorExit("Invalid mapping specification %s, from=to expected\n", mapping)
		}
		from := strings.TrimSpace(parts[0])
		to := strings.TrimSpace(parts[1])
		m[from] = to
	}
}

func readMapFile(filePath string, globalMappings map[string]string, fileMappings map[string]map[string]string) {
	file, err := os.Open(filePath)
	if err != nil {
		errorExit("Cannot open file \"%s\": %s\n", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	groupPattern := regexp.MustCompile(`\[(.+)]`)
	mappingPattern := regexp.MustCompile(`(.+?)\s*=\s*(.+)`)
	mappings := globalMappings
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		match := mappingPattern.FindStringSubmatch(line)
		if match != nil {
			mappings[match[1]] = match[2]
			continue
		}
		match = groupPattern.FindStringSubmatch(line)
		if match != nil {
			var ok bool
			if mappings, ok = fileMappings[match[1]]; !ok {
				mappings = make(map[string]string)
				fileMappings[match[1]] = mappings
			}
		}
	}

	if err := scanner.Err(); err != nil {
		errorExit("Error while reading mapping file: %s\n", err)
	}
}

func errorExit(format string, a ...any) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

type stringArrayFlag []string

func (i *stringArrayFlag) String() string {
	return "-"
}

func (i *stringArrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}
