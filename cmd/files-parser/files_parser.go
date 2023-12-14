package main

import (
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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.String("file_type", "", "Допустимые значения: csv")
	pflag.String("file_path", "", "")
	pflag.Int64("from_byte", 0, "")
	pflag.Int64("rows_limit", 1, "")
	pflag.Int64("threads_count", 1, "")
	pflag.String("mem_profile_path", "", "")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	fileType := viper.GetString("file_type")
	filePath := viper.GetString("file_path")
	fromByte := viper.GetInt64("from_byte")
	rowsLimit := viper.GetInt64("rows_limit")
	threadsCount := viper.GetInt64("threads_count")
	memProfilePath := viper.GetString("mem_profile_path")

	if fileType != "csv" {
		log.Fatalln("Тип файла не поддерживается")
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	batchesFromByte, err := csvparser.NewParser(file, fromByte, rowsLimit, threadsCount).CalculateBatchesFromByte()
	if err != nil {
		log.Fatalln(err)
	}

	defer ants.Release()
	var wg sync.WaitGroup

	for i := 0; i <= len(batchesFromByte)-1; i++ {
		i := i
		task := func() {
			file, err := os.Open(filePath)
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()

			log.Println(csvparser.NewParser(file, batchesFromByte[i], rowsLimit, 1).Parse())

			profile(memProfilePath)
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
