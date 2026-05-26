package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

func simWithOracle(filename string, caches ...Cache) error {

	f, err := os.OpenFile(filename, os.O_RDONLY, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	for entryIdx := 0; ; entryIdx++ {
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
		for _, cache := range caches {
			cache.GetOrInsert(id, objSize)
		}
	}

	// var s uint64 = 0
	// var maxSize uint32 = 0
	// var minSize uint32 = math.MaxUint32
	// for _, size := range sizes {
	// 	s += uint64(size)
	// 	maxSize = max(maxSize, size)
	// 	minSize = min(minSize, size)
	// }
	// average := float64(s) / float64(totalEntries)

	// // var variance float32 = 0
	// var varS float64 = 0
	// for _, size := range sizes {

	// 	var diff float64 = float64(size) - average
	// 	varS += math.Pow(diff, 2)
	// }
	// variance := varS / float64(totalEntries-1)

	// fmt.Printf("Stats for %s\n", filename)
	// fmt.Printf("Average: %f\n", average)
	// fmt.Printf("Total unique ids: %d\n", len(ids))
	// fmt.Printf("Total Unique bytes: %d\n", totalUniqueBytes)
	// fmt.Printf("Variance: %f\n", variance)
	// fmt.Printf("StdDev: %f\n", math.Sqrt(variance))
	// fmt.Printf("Max size: %d (%d)\n", maxSize, int(math.Log2(float64(maxSize))))
	// fmt.Printf("Min size: %d (%d)\n", minSize, int(math.Log2(float64(minSize))))

	fmt.Println()
	for _, cache := range caches {
		fmt.Printf("================ %s ================\n", cache.GetName())
		fmt.Printf("Byte Hits: %d (MiB)\n", cache.GetByteHits()/MiB)
		fmt.Printf("Hits: %0.2f (k)\n", float64(cache.GetHits())/1000)
	}

	return nil
}
