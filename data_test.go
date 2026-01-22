package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestData(t *testing.T) {
	f, err := os.OpenFile("data.txt", os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	for i := range 10_032 {
		f.Write(fmt.Appendf(nil, "%d,", i+1))
	}

	expectedSum := 10_032 * 10_033 / 2

	b := make([]byte, 10, 1024)
	newOffset, err := f.Seek(0, io.SeekStart)
	fmt.Println(newOffset)
	if err != nil {
		t.Fatal(err)
	}
	bf := bufio.NewReaderSize(f, 10)
	leftover := 0
	currSum := 0

	for {
		n, err := bf.Read(b[leftover:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println(currSum, expectedSum)
				return
			}
			t.Fatal(err)
		}
		total := leftover + n
		fmt.Println(leftover, n, total, string(b[:total]))
		leftover = 0
		for i := range total {
			if b[i] == ',' {
				num, err := strconv.Atoi(string(b[leftover:i]))
				if err != nil {
					log.Fatal(err)
				}
				currSum += num
				leftover = i + 1
			}
		}
		copy(b, b[leftover:total])
		leftover = total - leftover
	}

}
