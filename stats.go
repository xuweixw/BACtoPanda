package main

import "fmt"

type FastqStat struct {
	name     string
	numSeq   int64
	sumBases int64
}

func (fs *FastqStat) String() string {
	return fmt.Sprintf("%s\t%d\t%d\t", fs.name, fs.numSeq, fs.sumBases)
}
