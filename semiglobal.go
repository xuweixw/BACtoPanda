package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/io/seqio/fastq"
	"github.com/biogo/biogo/seq/linear"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// struct array for scores
//
//	Substitution matrix	置换矩阵
//		s(a, b) if a == b, s=1, else s = 0, gap (-1)
//	Initialization
//		F(i, 0) = 0
//		F(0, j)	= 0
//	Iteration F(i, j) = max(F(i, j-1)-d, F(i-1, j)-d, F(i-1, j-1) -d)
//
type scoreMatrix struct {
	data     []int
	row, col int
}

func NewScoreMatrix(seqA_len, seqB_len int) scoreMatrix {
	row := seqA_len + 1
	col := seqB_len + 1
	length := row * col
	s := make([]int, length, length)
	return scoreMatrix{s, row, col}
}

func (sm scoreMatrix) string() (str string) {
	for r := 0; r < sm.row; r++ {
		for c := 0; c < sm.col; c++ {
			str += fmt.Sprintf("%3d", sm.data[r*sm.col+c])
		}
		str += fmt.Sprintf("\n")
	}
	return
}

func (sm *scoreMatrix) setElement(row, col, value int) {
	location := row*sm.col + col
	sm.data[location] = value
}

func (sm *scoreMatrix) getElement(row, col int) int {
	location := row*sm.col + col
	return sm.data[location]
}
func (sm *scoreMatrix) lastCol() (l []int) {
	for r := 0; r < sm.row; r++ {
		loc := (r+1)*sm.col - 1
		l = append(l, sm.data[loc])
	}
	return
}
func (sm *scoreMatrix) lastRow() (r []int) {
	for l := 0; l < sm.col; l++ {
		r = append(r, sm.data[(sm.row-1)*sm.col+l])
	}
	return
}

type SemiGlobal struct {
	seqA string
	seqB string
	m    scoreMatrix
}

// NewSemiGlobal function.
// Sequence a is in the left of matrix.
// Seqeunce b is on the top of matrix.
func NewSemiGlobal(a, b string) *SemiGlobal {
	var scoring int
	score := NewScoreMatrix(len(a), len(b))
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
	return &SemiGlobal{seqA: a, seqB: b, m: score}
}

func (sm SemiGlobal) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("SeqA: %s\t%d bp\n", sm.seqA, len(sm.seqA)))
	s.WriteString(fmt.Sprintf("SeqB: %s\t%d bp\n", sm.seqB, len(sm.seqB)))
	s.WriteString(fmt.Sprintf("Scoring Matrix: \n%s\n", sm.m.string()))
	return s.String()
}

func (sm SemiGlobal) IsOverlap(min, edit int) (*string, bool) {
	var indexMaxValue, maxValue int
	for i, v := range sm.m.lastRow() { // find maximum value in last row
		if i == 0 || v > maxValue {
			indexMaxValue = i
			maxValue = v
		}
	}
	if maxValue > min && indexMaxValue-maxValue <= edit {
		subSeq := string(sm.seqB[maxValue:])
		return &subSeq, true
	}
	return nil, false
}

// create a func to alignment of two DNA strings

// IsOverlap
// semi-global alignment
//	a refers to BAC vector end sequences
//	b refers to a NGS read
// Notes: IsOverlap function has been convert into SemiGlobal.IsOverlap method, Dec 22, 2021.
//func IsOverlap(a, b string, min, edit int) bool {
//	var scoring int
//	score := NewScoreMatrix(len(a), len(b))
//	for i := 0; i < len(a); i++ {
//		for j := 0; j < len(b); j++ {
//			// 比对分
//			if a[i] == b[j] {
//				scoring = 1
//			} else {
//				scoring = -1
//			}
//			// 打分记录
//			threeMax := []int{
//				score.getElement(i, j) + scoring,
//				score.getElement(i+1, j) - 1,
//				score.getElement(i, j+1) - 1,
//			}
//			sort.Ints(threeMax)
//			score.setElement(i+1, j+1, threeMax[2])
//		}
//	}
//
//	var indexMaxvalue, maxValue int
//	for i, v := range score.lastRow() {
//		if i == 0 || v > maxValue {
//			indexMaxvalue = i
//			maxValue = v
//		}
//	}
//	if maxValue > min && indexMaxvalue-maxValue <= edit {
//		return true
//	}
//	return false
//}

// 计算反向互补序列
func reverseComplement(s string) (rc string) {
	base := map[byte]byte{
		'A': 'T',
		'T': 'A',
		'G': 'C',
		'C': 'G',
	}
	n := len(s)
	var bs []byte
	for i := n - 1; i >= 0; i-- {
		bs = append(bs, base[s[i]])
	}
	rc = string(bs)
	return
}

type PEReads struct {
	R1, R2  *linear.Seq
	overlap bool
}

func (p PEReads) isOverlapping(VES string, min, edit int) bool {
	R1Forward := p.R1.Seq.String()
	R1RevereComp := reverseComplement(R1Forward)

	R2Forward := p.R2.Seq.String()
	R2ReverseComp := reverseComplement(R2Forward)

	l := len(VES)
	seq := VES[l-18 : l]
	seed, _ := regexp.Compile(seq)
	var firstRound bool
	//var blankSlice = []string{}
	// 使用18bp短序列快速筛选
	if n := seed.FindStringSubmatch(R1Forward); len(n) != 0 {
		firstRound = true
	} else if n := seed.FindStringSubmatch(R1RevereComp); len(n) != 0 {
		firstRound = true
	} else if n := seed.FindStringSubmatch(R2Forward); len(n) != 0 {
		firstRound = true
	} else if n := seed.FindStringSubmatch(R2ReverseComp); len(n) != 0 {
		firstRound = true
	} else {
		return false
	}
	// 在第一轮筛选基础上，对通过的序列进行半全局比对
	if firstRound == true {
		for _, v := range []string{R1Forward, R1RevereComp, R2Forward, R2ReverseComp} {
			sm := NewSemiGlobal(VES, v)
			if subseq, bool := sm.IsOverlap(min, edit); bool {
				//fmt.Println(*subseq) /////////////////此处调用函数查找序列位置
				var BES string
				if len(*subseq) > 20 {
					BES = (*subseq)[6:]
					BWFM.update(FS, BES)
				}
				return bool
			}
		}
	}
	return false
}
func (p PEReads) Writer(writer *fasta.Writer) {
	if _, err := writer.Write(p.R1); err != nil {
		log.Fatal(err)
	}
	if _, err := writer.Write(p.R2); err != nil {
		log.Fatal(err)
	}
}

// ReadPEFastq function read paired-end reads into jobs channel.
// R1 and R2 arguments represent the path of R1 and R2 reads.
func ReadPEFastq(jobs chan<- *PEReads, R1, R2 string) {
	//defer wg.Done()
	t := linear.NewSeq("", nil, alphabet.DNA)

	R1_file, err := os.Open(R1)
	checkError(err)
	defer R1_file.Close()
	checkError(err)
	R2_file, err := os.Open(R2)
	checkError(err)
	defer R2_file.Close()
	uncompressR1_fq, err := gzip.NewReader(R1_file)
	checkError(err)
	uncompressR2_fq, err := gzip.NewReader(R2_file)
	checkError(err)
	R1_fq := fastq.NewReader(uncompressR1_fq, t)
	R2_fq := fastq.NewReader(uncompressR2_fq, t)

	R1_sc := seqio.NewScanner(R1_fq)
	R2_sc := seqio.NewScanner(R2_fq)
	var count int64
	for {
		R1_bool := R1_sc.Next()
		R2_bool := R2_sc.Next()
		if R1_bool == false || R2_bool == false {
			break
		}
		curr_PEReads := PEReads{R1: R1_sc.Seq().(*linear.Seq),
			R2:      R2_sc.Seq().(*linear.Seq),
			overlap: false}
		jobs <- &curr_PEReads
		count++
	}
	total <- count
	//return count
	//close(jobs)
}

/*
>pIndigoBAC536-S-forward
TCTTCGCTATTACGCCAGCTGGCGAAAGGGGGATGTGCTGCAAGGCGATTAAGTTGGGTAACGCCAGGGTTTTC
CCAGTCACGACGTTGTAAAACGACGGCCAGTGAATTGTAATACGACTCACTATAGGGCGAATTCGAGCTCGGTA
CCCGGGGATCCTCTAGAGTCGACCTGCAGGCATGC
*/

//var (
//	R1_file = flag.String("R1","", "NGS R1 reads file path")
//	R2_file = flag.String("R2","", "NGS R2 reads file path")
//	out_file = flag.String("out","", "target R1 and R2 has been merged in a fasta format file")
//	VES_seq = flag.String("ves","", "BAC vector end sequence, about 150 bp")
//	min_overlap = flag.Int("min", 20, "the min overlapping")
//	edit = flag.Int("edit", 1, "edit distance")
//	help = flag.Bool("help", false, "the help document")
//)

// The two vector end sequence of pIndigoBAC536-S are long 158 bp and 150 bp, respectively.
var (
	MIN_OVERLAP = 20
	EDIT        = 2
	jobChan     = make(chan *PEReads, 1000)
	resultChan  = make(chan *PEReads, 15)
	total       = make(chan int64, 1) //count reads number for stopping go concurrency
	t, n        int64
	wg          sync.WaitGroup
)

func worker(ctx context.Context, in <-chan *PEReads, out chan<- *PEReads, VES string, min, edit int) {
	//defer wg.Done()
	//blankPEReads := PEReads{nil, nil, false}
	for {
		select {
		case job := <-in:
			if job.isOverlapping(VES, min, edit) {
				job.overlap = true
				out <- job
			}
			atomic.AddInt64(&n, 1)
		case <-ctx.Done():
			return
		}
	}
	//results <- blank_PEReads
}

// FindOneVESBES function is generic compared to  FindTwoVESBES function.
func FindOneVESBES(R1, R2, VES, output string) {
	//flag.Parse()
	//if *help {
	//	flag.Usage()
	//}
	//wg.Add(1)
	t = 0 // initilization t and n value, in Go concurrency
	n = 0
	go ReadPEFastq(jobChan, R1, R2)
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < int(threads); i++ {
		//wg.Add(1)
		go worker(ctx, jobChan, resultChan, VES, MIN_OVERLAP, EDIT)
	}

	//wg.Add(1)
	//go func(o string) {
	//defer wg.Done()
	fmt.Println(output, VES)
	out, err := os.Create(output)
	checkError(err)
	defer out.Close()
	out_fasta := fasta.NewWriter(out, 150)
	for {
		select {
		case r := <-resultChan:
			if r.overlap {
				r.Writer(out_fasta)
			}
		case t = <-total:
			fmt.Println("Read over 1")
		case <-time.Tick(2 * time.Second):
			//fmt.Printf("%d o-------------k%d\n", t, n)
			if t != 0 && t == n {
				cancel()
				time.Sleep(5 * time.Second)
				return
			}
		}
	}
	//}(output)
	//wg.Wait()
}
