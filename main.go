package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
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

func main() {

	simWithOracle("meta_reag.oracleGeneral")

}
