package main

import "testing"

func TestCorrectBED(t *testing.T) {
	bbs := ReadBEDBaseSet("depthCalling/SP01-01_BAC.bed")
	//
	bwfs := ReadBEDWithFrequencySet("data_test/sub.bed")
	//
	bbs.Correct(bwfs, 2000, 5)
	//
	bbs.Write("Allcorrect.bed")
}
