package main

import (
	"testing"
)

func TestInsert(t *testing.T) {
	s := NewDLL()
	s.insert("1")
	s.insert("2")
	s.insert("3")
	s.insert("4")
}
func TestMoveToHead(t *testing.T) {
	s := NewDLL()
	m := s.insert("1")
	s.insert("2")
	s.insert("3")
	s.insert("4")
	s.moveToHead(m)
}
func TestRemove(t *testing.T) {
	s := NewDLL()
	s.insert("1")
	s.insert("2")
	s.insert("3")
	m := s.insert("4")
	s.remove(m)

}
func TestMoveAndRemove(t *testing.T) {
	s := NewDLL()
	m := s.insert("4")
	s.moveToHead(m)
	s.remove(m)

}
