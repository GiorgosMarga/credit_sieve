package main

// import (
// 	"errors"
// 	"fmt"
// 	"math"
// )

// type MyCache struct {
// 	internalCaches []Sieve
// 	size           int
// 	cap            int
// 	hits           int
// 	byteHit        int
// 	misses         int
// 	byteMissed     int
// }

// func NewMyCache(cap int) MyCache {
// 	// 1024 2^i-1 <= 1024
// 	totalCaches := int(math.Floor(math.Log2(float64(cap))))
// 	// totalCaches := cap
// 	fmt.Println(totalCaches)
// 	caches := make([]Sieve, 0, totalCaches)
// 	for range totalCaches {
// 		// from := int(math.Pow(float64(2), float64(i)))
// 		// to := int(math.Pow(float64(2), float64(i+1)) - 1)
// 		// fmt.Printf("%d -> %d-%d\n", i+1, from, to)
// 		caches = append(caches, NewSieve(cap*10, SieveOpts{}))
// 	}
// 	return MyCache{
// 		internalCaches: caches,
// 		cap:            cap,
// 	}
// }

// func (c *MyCache) GetOrInsert(k string) (any, error) {
// 	size := len(k)
// 	if size == 0 {
// 		return nil, nil
// 	}
// 	cacheId := int(math.Floor(math.Log2(float64(size))))
// 	res, err := c.internalCaches[cacheId].Get(k)
// 	if err != nil {
// 		if !errors.Is(err, ErrNotFound) {
// 			return nil, err
// 		}
// 	}
// 	if res != nil {
// 		c.hits++
// 		c.byteHit += len(k)
// 		return res, nil
// 	}
// 	c.misses++
// 	c.byteMissed += len(k)

// 	var remId int = 0
// 	if c.size+len(k)-c.cap > 0 {
// 		remId = int(math.Floor(math.Log2(float64(c.size + len(k) - c.cap))))
// 	}

// 	for c.size+len(k) > c.cap {
// 		evictedSize := c.internalCaches[remId].Evict()
// 		if evictedSize == 0 {
// 			remId = (remId - 1 + len(c.internalCaches)) % len(c.internalCaches)
// 			// remId = (remId + 1) % len(c.internalCaches)
// 			continue
// 		}
// 		c.size -= evictedSize
// 		remId = int(math.Floor(math.Log2(float64(c.size + len(k) - c.cap))))
// 	}

// 	if err := c.internalCaches[cacheId].Put(k, ""); err != nil {
// 		return nil, err
// 	}
// 	c.size += len(k)

// 	return "", nil
// }
