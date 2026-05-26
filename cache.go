package main

type Cache interface {
	GetOrInsert(key string, size uint32) (any, error)
	GetByteHits() uint32
	GetHits() uint32
	GetName() string
}
