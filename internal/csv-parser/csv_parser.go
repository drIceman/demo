package csvparser

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"os"
)

func Parse(filePath string, fromByte int64, rowsLimit int64) (rows [][]string, newFromByte int64, isEOF bool, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return rows, fromByte, isEOF, err
	}
	defer file.Close()

	newFromByte, rowsLimit, err = calculateCorrectFromByteAndRowsLimit(file, fromByte, rowsLimit)
	if err != nil {
		return rows, fromByte, isEOF, err
	}
	if rowsLimit == 0 {
		return rows, newFromByte, true, err
	}

	reader := csv.NewReader(bufio.NewReader(file))
	rows = make([][]string, rowsLimit)
	for i := int64(0); i < rowsLimit; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			isEOF = true
			break
		}
		if err != nil {
			return rows, rowsLimit, isEOF, err
		}

		rows[i] = row
	}

	return rows, newFromByte + reader.InputOffset(), isEOF, err
}

func calculateCorrectFromByteAndRowsLimit(file *os.File, fromByte int64, rowsLimit int64) (newFromByte int64, newRowsLimit int64, err error) {
	fInfo, err := file.Stat()
	if err != nil {
		return fromByte, rowsLimit, err
	}

	fSize := fInfo.Size()
	if fromByte >= fSize {
		return fSize, 0, err
	}

	_, err = file.Seek(fromByte, 0)
	if err != nil {
		return fromByte, rowsLimit, err
	}

	bufSize := bufio.MaxScanTokenSize
	if diff := fSize - fromByte; diff < int64(bufSize) {
		bufSize = int(diff)
	}

	buf := make([]byte, bufSize)
	for {
		bufSize, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fromByte, rowsLimit, err
		}

		var bufPosition int
		for {
			i := bytes.IndexByte(buf[bufPosition:], '\n')
			if i == -1 || bufSize == bufPosition || newRowsLimit >= rowsLimit {
				break
			}
			bufPosition += i + 1
			newRowsLimit++
		}
	}

	newFromByte, err = file.Seek(fromByte, 0)
	if err != nil {
		return fromByte, rowsLimit, err
	}

	return newFromByte, newRowsLimit, err
}
