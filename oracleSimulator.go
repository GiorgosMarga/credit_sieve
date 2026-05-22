package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

func simWithOracle(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	totalBytes := stat.Size()
	totalEntries := totalBytes / 24
	printEvery := totalEntries/10 + 1

	sieve := NewSieve(2 * MiB)
	gmr := NewGMR(2 * MiB)

	sizes := make([]uint32, totalEntries)

	for entryIdx := 0; ; entryIdx++ {
		if entryIdx%int(printEvery) == 0 {
			percent := (entryIdx * 100) / int(totalEntries)
			fmt.Printf("%d%% done (%d / %d entries)\n", percent, entryIdx, totalEntries)
		}

		buf := make([]byte, 24)
		n, err := f.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if n != 24 {
			panic(n)
		}
		objId := binary.LittleEndian.Uint64(buf[4:])
		objSize := binary.LittleEndian.Uint32(buf[12:])

		id := fmt.Sprintf("%d", objId)

		sieve.GetOrInsert(id, objSize)
		gmr.GetOrInsert(id, objSize)
		sizes[entryIdx] = objSize
	}

	var s uint64 = 0
	var maxSize uint32 = 0
	var minSize uint32 = math.MaxUint32
	for _, size := range sizes {
		s += uint64(size)
		maxSize = max(maxSize, size)
		minSize = min(minSize, size)
	}
	average := float64(s) / float64(totalEntries)

	// var variance float32 = 0
	var varS float64 = 0
	for _, size := range sizes {

		var diff float64 = float64(size) - average
		varS += math.Pow(diff, 2)
	}
	variance := varS / float64(totalEntries-1)

	fmt.Printf("Stats for %s\n", filename)
	fmt.Printf("Average: %f\n", average)
	fmt.Printf("Variance: %f\n", variance)
	fmt.Printf("StdDev: %f\n", math.Sqrt(variance))
	fmt.Printf("Max size: %d (%d)\n", maxSize, int(math.Log2(float64(maxSize))))
	fmt.Printf("Min size: %d (%d)\n", minSize, int(math.Log2(float64(minSize))))

	fmt.Printf("Byte Hits: GMR %d vs %d SIEVE (MiB)\n", gmr.byteHits/MiB, sieve.byteHits/MiB)
	fmt.Printf("Hits: GMR %d vs %d SIEVE\n", gmr.hits, sieve.hits)
	return nil
}
