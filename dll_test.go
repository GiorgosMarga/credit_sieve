package main

import (
	"testing"
)

func TestInsert(t *testing.T) {
	s := NewDLL()
	s.insert("1", 1)
	s.insert("2", 1)
	s.insert("3", 1)
	s.insert("4", 1)
}
func TestMoveToHead(t *testing.T) {
	s := NewDLL()
	m := s.insert("1", 1)
	s.insert("2", 1)
	s.insert("3", 1)
	s.insert("4", 1)
	s.moveToHead(m)
}
func TestRemove(t *testing.T) {
	s := NewDLL()
	s.insert("1", 1)
	s.insert("2", 1)
	s.insert("3", 1)
	m := s.insert("4", 1)
	s.remove(m)

}
func TestMoveAndRemove(t *testing.T) {
	s := NewDLL()
	m := s.insert("4", 1)
	s.moveToHead(m)
	s.remove(m)

}
