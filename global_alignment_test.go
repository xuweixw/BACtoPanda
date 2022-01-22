package main

import (
	"fmt"
	"testing"
)

func TestNewGlobalAligner(t *testing.T) {
	ga := NewGlobalAligner("ATAGTATGCATTCATAAGCTATC", "ATAGTATGCATTCATGCTATAC")
	fmt.Println(ga, ga.IsSimilar(2))
}
