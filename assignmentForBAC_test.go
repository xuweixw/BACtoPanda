package main

import (
	"fmt"
	"testing"
)

func TestBorder(t *testing.T) {
	n := 150
	for i := 1; i <= 100; i++ {
		margin = Margin(n * i)
		_, n := Location()
		fmt.Printf("%.0f\t%d\n", margin, n)
	}
}

func TestIntersection(t *testing.T) {
	//Hic_asm_0	12555930	12665226
	//Hic_asm_0	20211144	20317010
	//Hic_asm_0	22495321	22609245
	//Hic_asm_0	24960826	25077423
	//Hic_asm_0	41105250	41209022
	r1 := Region{chr: "Hic_asm_1", start: 12555930, end: 12665226}
	r2 := Region{chr: "Hic_asm_0", start: 20211144, end: 20317010}
	r3 := Region{chr: "Hic_asm_1", start: 22495321, end: 22609245}
	r4 := Region{chr: "Hic_asm_0", start: 24960826, end: 25077423}
	r5 := Region{chr: "Hic_asm_0", start: 41105250, end: 41209022}
	set1 := RegionSet{regions: []Region{r1, r2, r3}}
	set2 := RegionSet{regions: []Region{r3, r4, r5, r1}}
	fmt.Printf("%v", Intersection(&set1, &set2))
}

func TestAssignment(t *testing.T) {
	destinationDir = "."
	//plateID = [7]string{"8", "9","10","11","12","13","14"}
	//plateID = [7]int{1,2,3,4,5,6,7}
	//if err:=  os.Mkdir( "assignment", 0700); err != nil{
	//	t.Errorf("create error: %v", err)
	//}
	Assignment()
}
func TestLocation(t *testing.T) {
	destinationDir = "."

	for i := 100; i < 5000; i += 50 {
		margin = Margin(float64(i))
		r, num := Location()
		fmt.Printf("%d Get, %d Lose in border=%d\n", len(r.regions), num, i)
	}
}
