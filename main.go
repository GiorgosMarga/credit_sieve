package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	_ = 1 << (10 * iota) // ignore first value

	KiB // 1 KiB = 1024 bytes
	MiB // 1 MiB = 1024 KiB
	GiB // 1 GiB = 1024 MiB
	TiB // 1 TiB = 1024 GiB
	PiB
)

func lazyCreateTestFile(numLines int) string {
	fn := fmt.Sprintf("numbers.%d.txt", numLines)
	if _, err := os.Stat(fn); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return generateFile(fn, numLines)
		}
		return ""
	}
	return fn
}

func generateFile(s string, numLines int) string {
	f, err := os.OpenFile(s, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// s = 0.8, v = 1, imax = N-1
	z := rand.NewZipf(r, 1.1, 3, math.MaxUint64)
	for i := range numLines {
		if i == numLines-1 {
			f.Write(fmt.Appendf(nil, "%d", z.Uint64()))
		} else {
			f.Write(fmt.Appendf(nil, "%d,", z.Uint64()))
		}
	}
	return s
}

var percentages []float64 = []float64{0.1, 0.5, 1, 5, 10, 20}

func main() {
	_ = percentages
	var (
		filename  string
		cacheSize uint32 = 2 * GiB
	)

	flag.StringVar(&filename, "f", "meta_reag.oracleGeneral", "Cache file to run.")
	flag.Func("c", "Cache size (kb | mb | gb). Ex. 2gb -> 2 gigabytes.", func(s string) error {
		if s == "" {
			panic("invalid cache size")
		}
		cacheSizeStr := s[:len(s)-2]
		size, err := strconv.Atoi(cacheSizeStr)
		if err != nil {
			panic(err)
		}
		scale := s[len(s)-2:]
		switch scale {
		case "kb":
			cacheSize = uint32(size) * KiB
		case "mb":
			cacheSize = uint32(size) * MiB
		case "gb":
			cacheSize = uint32(size) * GiB
		default:
			panic("invalid scale")
		}
		fmt.Println(cacheSize)
		return nil
	})
	flag.Parse()

	totalBytes, err := getTotalBytes(filename)
	if err != nil {
		log.Fatal(err)
	}

	caches := make([]Cache, 0)
	for _, percentage := range percentages {
		c := uint32(float64(totalBytes) * percentage / 100)
		caches = append(caches, NewLRU(c))
		caches = append(caches, NewGMR(c))
		caches = append(caches, NewSieve(c))

	}

	simWithOracle(filename, caches...)

}

func getTotalBytes(filename string) (uint, error) {
	f, err := os.OpenFile("stats.csv", os.O_RDONLY, 0o666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	_, err = csvReader.Read()
	if err != nil {
		return 0, err
	}
	for {
		row, err := csvReader.Read()
		if err != nil {
			return 0, err
		}
		name := row[0]
		if name == strings.Split(filename, ".")[0] {
			totalBytesStr := row[4]
			t, err := strconv.ParseUint(totalBytesStr, 10, 64)
			if err != nil {
				return 0, err
			}
			return uint(t), nil
		}
	}
}
