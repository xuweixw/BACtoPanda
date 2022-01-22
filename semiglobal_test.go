package main

import (
	"fmt"
	"testing"
)

//func TestIsOverlap(t *testing.T) {
//	r := IsOverlap("AAATCGATCGATCAAATCGAT", "CAAATCGATAAATCGATACATTATCGAT", 2,0)
//	fmt.Println("result: ", r)
//}

func TestNewSemiGlobal(t *testing.T) {
	sm := NewSemiGlobal("ATAGCTATATTCGTAC", "CTATATTCGTACATCGATTCGATCTA")
	sub, ok := sm.IsOverlap(4, 0)
	fmt.Println(*sub, ok, sm)
}
