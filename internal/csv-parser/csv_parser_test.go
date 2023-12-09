package csvparser

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		filePath  string
		fromByte  int64
		rowsLimit int64
	}

	var emptyRecords [][]string
	tests := []struct {
		name          string
		args          args
		wantRecords   [][]string
		wantNewOffset int64
		wantIsEOF     bool
		wantErr       bool
	}{
		{
			name:          "Парсинг одной строки без сдвига",
			args:          args{filePath: "./stub.csv", fromByte: 0, rowsLimit: 1},
			wantRecords:   [][]string{{"Поле 1", "Поле 2", "Поле 3", "Поле 4", "Поле 5", "Поле 6", "Поле 7", "Поле 8", "Поле 9", "Поле 10"}},
			wantNewOffset: 131,
			wantIsEOF:     false,
			wantErr:       false,
		},
		{
			name: "Парсинг двух строк без сдвига",
			args: args{filePath: "./stub.csv", fromByte: 0, rowsLimit: 2},
			wantRecords: [][]string{
				{"Поле 1", "Поле 2", "Поле 3", "Поле 4", "Поле 5", "Поле 6", "Поле 7", "Поле 8", "Поле 9", "Поле 10"},
				{"Поле 11", "Поле 12", "Поле 13", "Поле 14", "Поле 15", "Поле 16", "Поле 17", "Поле 18", "Поле 19", "Поле 20"},
			},
			wantNewOffset: 271,
			wantIsEOF:     false,
			wantErr:       false,
		},
		{
			name:          "Парсинг одной строки со сдвигом",
			args:          args{filePath: "./stub.csv", fromByte: 1392, rowsLimit: 1},
			wantRecords:   [][]string{{"Поле 101", "Поле 102", "Поле 103", "Поле 104", "Поле 105", "Поле 106", "Поле 107", "Поле 108", "Поле 109", "Поле 110"}},
			wantNewOffset: 1542,
			wantIsEOF:     false,
			wantErr:       false,
		},
		{
			name:          "Парсинг в конце с превышением числа строк",
			args:          args{filePath: "./stub.csv", fromByte: 1542, rowsLimit: 2},
			wantRecords:   [][]string{nil},
			wantNewOffset: 1543,
			wantIsEOF:     true,
			wantErr:       false,
		},
		{
			name:          "Парсинг в конце с превышением сдвига",
			args:          args{filePath: "./stub.csv", fromByte: 1545, rowsLimit: 1},
			wantRecords:   emptyRecords,
			wantNewOffset: 1543,
			wantIsEOF:     true,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRecords, gotNewOffset, gotIsEOF, err := Parse(tt.args.filePath, tt.args.fromByte, tt.args.rowsLimit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRecords, tt.wantRecords) {
				t.Errorf("Parse() gotRecords = %v, want %v", gotRecords, tt.wantRecords)
			}
			if gotNewOffset != tt.wantNewOffset {
				t.Errorf("Parse() gotNewOffset = %v, want %v", gotNewOffset, tt.wantNewOffset)
			}
			if gotIsEOF != tt.wantIsEOF {
				t.Errorf("Parse() gotIsEOF = %v, want %v", gotIsEOF, tt.wantIsEOF)
			}
		})
	}
}

func TestPrepare1G(t *testing.T) {
	if os.Getenv("TEST_1G") != "1" {
		t.SkipNow()
	}
	fIn, _ := os.Open("./stub.csv")
	fOut, _ := os.Create("./stub_1g.csv")
	content, _ := io.ReadAll(fIn)
	for i := 0; i < 1000000; i++ {
		io.WriteString(fOut, string(content))
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("./stub.csv", 0, 1)
	}
}

func BenchmarkParse1G(b *testing.B) {
	if os.Getenv("TEST_1G") != "1" {
		b.SkipNow()
	}
	b.StartTimer()
	fromByte := int64(0)
	for {
		_, newFromByte, isEOF, err := Parse("./stub_1g.csv", fromByte, 1000)
		if err != nil || isEOF == true {
			break
		}
		fromByte = newFromByte
		b.Log("Текущая позиция: ", fromByte)
	}
	b.StopTimer()
	b.Log("Времени затрачено (мс): ", b.Elapsed().Microseconds())
}
