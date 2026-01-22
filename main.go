package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
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

func Evict(s *Sieve) int {
	if len(s.items) == 0 {
		return 0
	}
	o := s.hand
	//tail is a dummy node
	if o == s.dll.tail || o == s.dll.head {
		o = s.dll.tail.prev

	}
	for o.visited == true {
		o.visited = false
		o = o.prev
		// head is a dummy node
		// wrap around list
		if o == s.dll.head {
			o = s.dll.tail.prev
		}
	}
	s.hand = o.prev
	s.dll.remove(o) // eviction
	s.size -= len(o.key)
	delete(s.items, o.val)
	return len(o.key)
}

func EvictByCredit(s *Sieve) int {
	if len(s.items) == 0 {
		return 0
	}
	o := s.hand
	//tail is a dummy node
	if o == s.dll.tail || o == s.dll.head {
		o = s.dll.tail.prev

	}
	for o.credit > 0 {
		// fmt.Println(o.credit)
		o.credit -= o.penalty
		o.penalty *= 2
		o = o.prev
		// head is a dummy node
		// wrap around list
		if o == s.dll.head {
			o = s.dll.tail.prev
		}
	}
	s.hand = o.prev
	s.dll.remove(o) // eviction
	s.size -= len(o.key)
	delete(s.items, o.val)
	return len(o.key)
}

func main() {
	// s := NewMyCache(1024 * 10)
	s := NewSieve(1024*10, SieveOpts{
		Evict: Evict,
	})
	s2 := NewSieve(1024*10, SieveOpts{
		Evict: EvictByCredit,
	})
	filename := lazyCreateTestFile(2_000_000)
	// s := NewSieve(1024)
	f, err := os.OpenFile(filename, os.O_RDONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b := make([]byte, 1024)
	numbersRead := 0
	bf := bufio.NewReaderSize(f, 1024)
	leftover := 0
	for {
		n, err := bf.Read(b[leftover:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if leftover != 0 {
					numbersRead++
					s.GetOrInsert(string(b[:leftover]))
					s2.GetOrInsert(string(b[:leftover]))
				}
				break
			}
			log.Fatal(err)
		}
		total := leftover + n
		leftover = 0
		for i := range total {
			if b[i] == ',' {
				s.GetOrInsert(string(b[leftover:i]))
				s2.GetOrInsert(string(b[leftover:i]))
				numbersRead++
				leftover = i + 1
			}
		}
		copy(b, b[leftover:total])
		leftover = total - leftover
	}
	fmt.Printf("Credit vs Original: %d vs %d\n", s2.byteHits, s.byteHits)
	fmt.Printf("Credit vs Original (Misses): %d vs %d\n", s2.byteMisses, s.byteMisses)

}
