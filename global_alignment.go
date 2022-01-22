package main

import (
	"fmt"
	"sort"
	"strings"
)

type GlobalAligner struct {
	seqA, seqB string
	m          scoreMatrix
}

func NewGlobalAligner(a, b string) *GlobalAligner {
	var scoring int
	score := NewScoreMatrix(len(a), len(b))
	for i := 0; i < len(a)+1; i++ { // initialize first column
		score.setElement(i, 0, -i)
	}
	for i := 0; i < len(b)+1; i++ { // initialize fisrt row
		score.setElement(0, i, -i)
	}

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			// 比对分
			if a[i] == b[j] {
				scoring = 1
			} else {
				scoring = -1
			}
			// 打分记录
			threeMax := []int{
				score.getElement(i, j) + scoring,
				score.getElement(i+1, j) - 1,
				score.getElement(i, j+1) - 1,
			}
			sort.Ints(threeMax)
			score.setElement(i+1, j+1, threeMax[2])
		}
	}
	return &GlobalAligner{seqA: a, seqB: b, m: score}
}

func (ga *GlobalAligner) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("SeqA: %s\t%d bp\n", ga.seqA, len(ga.seqA)))
	s.WriteString(fmt.Sprintf("SeqB: %s\t%d bp\n", ga.seqB, len(ga.seqB)))
	s.WriteString(fmt.Sprintf("Scoring Matrix: \n%s\n", ga.m.string()))
	return s.String()
}

func (ga *GlobalAligner) IsSimilar(edit int) bool {
	// find maximum score in last row
	var max int
	for _, v := range ga.m.lastRow() {
		if v > max {
			max = v
		}
	}
	if max+2*edit < len(ga.seqB) { // mismatch punish one score, so edit multiply 2.
		return false
	}
	return true
}

// Determine whether two sequence is equal, brutally, simply and efficiently.
func equal(a, b string) bool { // len(a) == len(b)
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
