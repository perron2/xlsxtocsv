xlsxtocsv
=========

Command line tool to convert Excel .xlsx files to CSV files.

Build
-----
With Go installed execute

```
go build -ldflags="-s -w" -trimpath perron2.ch/xlsxtocsv
```

You can also use the associated makefile if you have make
and [UPX](https://upx.github.io) installed on your system.
MacOS builds do currently not use UPX to compress the generated
executable file due to an unresolved problem in UPX on macOS 13+.

```
make windows
make linux
make darwin
```

Examples
--------

Convert a single Excel file to CSV and output the result to STDOUT:

```sh
xlsxtocsv sample.xlsx
```

Convert a single Excel file to CSV and output the result to a specific file:

```sh
xlsxtocsv -out=sample.csv sample.xlsx
```

Convert all .xlsx files in a directory at once and store the CSV files in
another directory:

```sh
xlsxtocsv -outdir=csv *.xlsx
```

Use ANSI (Windows-1252) character encoding instead of UTF-8 (default).
In addition, use the semicolon character instead of a comma as a field separator:

```sh
xlsxtocsv -separator=; -charset=ansi -out=sample.csv sample.xlsx
```

Rename columns "firstname" and "city" to "name" and "place":

```sh
xlsxtocsv -map firstname=name -map city=place sample.xlsx
```

Rename columns when converting multiple files at once using a separate
mapping file:

```sh
xlsxtocsv -mapfile=mappings.xls -outdir=csv *.xlsx
```

The mapping files might look like this:

```
[sample.xlsx]
firstname=name
city=place

[sample2.xlsx]
referenceNumber=id
```
