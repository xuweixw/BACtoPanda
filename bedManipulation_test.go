package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestWhat(t *testing.T) {
	file, err := os.Open("./test-data/BAC-40-1.depth.tab")
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	output := NewTiles(file)

	for _, item := range output.Splite(1000) {
		if item.potentialBAC() {
			fmt.Printf("%v", item.ToBed())
		}
	}
}

func TestCommand(t *testing.T) {
	var (
		input = flag.String("input", "", "depth.tab file")
		help  = flag.Bool("help", false, "help usage")
	)
	flag.Parse()
	if *help == true {
		flag.Usage()
		os.Exit(0)
	}
	file, err := os.Open(*input)
	defer file.Close()
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	output := NewTiles(file)

	for _, item := range output.Splite(1000) {
		if item.potentialBAC() {
			fmt.Printf("%v", item.ToBed())
		}
	}
}

func TestProximityInt(t *testing.T) {
	//sequence1 := []int{4,70,1,5,2,8,24,17}
	sequence2 := []int{}
	point := 18
	fmt.Println(ProximityInt(point, sequence2))
}

func TestCorrect(t *testing.T) {
	bwf1 := BEDWithFrequency{
		BEDsp{
			BEDBase{chr: "Chr1", start: 2332, end: 0}, ""}, 20}
	bwf2 := BEDWithFrequency{
		BEDsp{
			BEDBase{chr: "Chr1", start: 4345, end: 0}, ""}, 20}
	bwf3 := BEDWithFrequency{
		BEDsp{
			BEDBase{chr: "Chr2", start: 2324, end: 0}, ""}, 20}
	bwfs := BEDWithFrequencySet{ele: []BEDWithFrequency{bwf1, bwf2, bwf3}}
	bb1 := BEDBase{"Chr1", 2333, 2344}
	bb2 := BEDBase{"Chr2", 3333, 2300}
	bb3 := BEDBase{"Chr1", 4333, 2330}
	bbs := BEDBaseSet{[]*BEDBase{&bb1, &bb2, &bb3}}
	fmt.Println(bbs.String())
	//校准边界
	bbs.Correct(&bwfs, 200, 10)
	fmt.Println(bbs.String())
	x := ReadBEDWithFrequencySet("data_test/sub.bed")
	fmt.Println(x.String())
}
