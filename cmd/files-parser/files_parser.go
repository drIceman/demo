package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	csvparser "github.com/drIceman/demo/internal/csv-parser"
)

var fileType = flag.String("fileType", "", "Допустимые значения: csv")
var filePath = flag.String("filePath", "", "")
var fromByte = flag.Int64("fromByte", 0, "")
var rowsLimit = flag.Int64("rowsLimit", 0, "")
var memProfilePath = flag.String("memProfilePath", "", "")

func main() {
	flag.Parse()

	switch *fileType {
	case "csv":
		fmt.Println(csvparser.Parse(*filePath, *fromByte, *rowsLimit))
	default:
		log.Fatal("Тип файла не поддерживается")
	}

	if *memProfilePath != "" {
		f, err := os.Create(*memProfilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		runtime.GC()
		//pprof.Lookup("heap").WriteTo(f, 2)
		pprof.Lookup("allocs").WriteTo(f, 2)
		return
	}
}
