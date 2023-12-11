package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	csvparser "github.com/drIceman/demo/internal/csv-parser"
	"github.com/panjf2000/ants/v2"
)

func main() {
	fileType := flag.String("fileType", "", "Допустимые значения: csv")
	filePath := flag.String("filePath", "", "")
	fromByte := flag.Int64("fromByte", 0, "")
	rowsLimit := flag.Int64("rowsLimit", 1, "")
	threadsCount := flag.Int64("threadsCount", 1, "")
	memProfilePath := flag.String("memProfilePath", "", "")

	flag.Parse()

	// *fileType = "csv"
	// *filePath = "../../internal/csv-parser/stub.csv"
	// *fromByte = 0
	// *rowsLimit = 2
	// *threadsCount = 2
	// *memProfilePath = ""

	//stub
	//+ -fromByte=0 -rowsLimit=1 -threadsCount=5
	//+ -fromByte=0 -rowsLimit=5 -threadsCount=1
	//+ -fromByte=131 -rowsLimit=5 -threadsCount=1
	//+- -fromByte=131 -rowsLimit=1 -threadsCount=5
	//+ -fromByte=131 -rowsLimit=2 -threadsCount=2
	//+ -fromByte=0 -rowsLimit=2 -threadsCount=2

	if *fileType != "csv" {
		log.Fatalln("Тип файла не поддерживается")
	}

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	batchesFromByte, err := csvparser.CalculateBatchesFromByte(file, *fromByte, *rowsLimit, *threadsCount)
	if err != nil {
		log.Fatalln(err)
	}

	defer ants.Release()
	var wg sync.WaitGroup
	for _, fromByte := range batchesFromByte {
		fromByte := fromByte
		task := func() {
			log.Println(csvparser.Parse(file, fromByte, *rowsLimit))
			profile(*memProfilePath)
			wg.Done()
		}
		wg.Add(1)
		_ = ants.Submit(task)
	}
	wg.Wait()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Завершаем работу...")
	time.Sleep(1 * time.Second)
}

func profile(memProfilePath string) {
	if memProfilePath != "" {
		f, err := os.Create(memProfilePath)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		runtime.GC()
		//pprof.Lookup("heap").WriteTo(f, 2)
		pprof.Lookup("allocs").WriteTo(f, 2)
		os.Exit(0)
	}
}
