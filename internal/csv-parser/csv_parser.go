package csvparser

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"os"
)

type parser struct {
	file         *os.File
	fromByte     int64
	rowsLimit    int64
	threadsCount int64
}

func NewParser(file *os.File, fromByte int64, rowsLimit int64, threadsCount int64) parser {
	return parser{
		file:         file,
		fromByte:     fromByte,
		rowsLimit:    rowsLimit,
		threadsCount: threadsCount,
	}
}

func (p parser) Parse() (rows [][]string, newFromByte int64, isEOF bool, err error) {
	var newRowsLimit int64
	newFromByte, newRowsLimit, err = p.calculateCorrectFromByteAndRowsLimit()
	if err != nil {
		return rows, p.fromByte, isEOF, err
	}
	if newRowsLimit == 0 {
		return rows, newFromByte, true, err
	}

	reader := csv.NewReader(bufio.NewReader(p.file))
	rows = make([][]string, newRowsLimit)
	for i := int64(0); i < newRowsLimit; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			isEOF = true
			break
		}
		if err != nil {
			return rows, p.rowsLimit, isEOF, err
		}

		rows[i] = row
	}

	return rows, newFromByte + reader.InputOffset(), isEOF, err
}

func (p parser) CalculateBatchesFromByte() (batchesFromByte []int64, err error) {
	fInfo, err := p.file.Stat()
	if err != nil {
		return batchesFromByte, err
	}

	fSize := fInfo.Size()
	if p.fromByte >= fSize {
		return batchesFromByte, err
	}

	_, err = p.file.Seek(p.fromByte, 0)
	if err != nil {
		return batchesFromByte, err
	}

	bufSize := bufio.MaxScanTokenSize
	if diff := fSize - p.fromByte; diff < int64(bufSize) {
		bufSize = int(diff)
	}

	batchesFromByte = make([]int64, p.threadsCount)
	batchesFromByte[0] = p.fromByte
	buf := make([]byte, bufSize)
	newRowsLimit := int64(0)
	for i := int64(1); i < p.threadsCount; {
		bufSize, err := p.file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return batchesFromByte, err
		}

		var bufPosition int
		for {
			j := bytes.IndexByte(buf[bufPosition:], '\n')
			if j == -1 || bufSize == bufPosition || i == p.threadsCount {
				p.fromByte += int64(bufSize - bufPosition)
				break
			}
			bufPosition += j + 1
			p.fromByte += int64(j + 1)
			newRowsLimit++
			if newRowsLimit == p.rowsLimit {
				batchesFromByte[i] = p.fromByte
				newRowsLimit = 0
				i++
			}
		}
	}

	_, err = p.file.Seek(p.fromByte, 0)
	if err != nil {
		return batchesFromByte, err
	}

	return batchesFromByte, err
}

func (p parser) calculateCorrectFromByteAndRowsLimit() (newFromByte int64, newRowsLimit int64, err error) {
	fInfo, err := p.file.Stat()
	if err != nil {
		return p.fromByte, p.rowsLimit, err
	}

	fSize := fInfo.Size()
	if p.fromByte >= fSize {
		return fSize, 0, err
	}

	_, err = p.file.Seek(p.fromByte, 0)
	if err != nil {
		return p.fromByte, p.rowsLimit, err
	}

	bufSize := bufio.MaxScanTokenSize
	if diff := fSize - p.fromByte; diff < int64(bufSize) {
		bufSize = int(diff)
	}

	buf := make([]byte, bufSize)
	for {
		bufSize, err := p.file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return p.fromByte, p.rowsLimit, err
		}

		var bufPosition int
		for {
			i := bytes.IndexByte(buf[bufPosition:], '\n')
			if i == -1 || bufSize == bufPosition || newRowsLimit == p.rowsLimit {
				break
			}
			bufPosition += i + 1
			newRowsLimit++
		}
	}

	newFromByte, err = p.file.Seek(p.fromByte, 0)
	if err != nil {
		return p.fromByte, p.rowsLimit, err
	}

	return newFromByte, newRowsLimit, err
}
