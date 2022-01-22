package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type BEDBase struct { // Basic BED format (Browser Extensible Data), only fist three columns.
	chr        string
	start, end int
}

func (bb *BEDBase) String() string {
	return fmt.Sprintf("%s\t%d\t%d", bb.chr, bb.start, bb.end)
}

type BEDBaseSet struct {
	record []*BEDBase
}

func ReadBEDBaseSet(filePath string) *BEDBaseSet {
	var set = BEDBaseSet{record: make([]*BEDBase, 0)}
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		log.Fatalln("Open file Error: ", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := strings.Split(scanner.Text(), "\t")
		chr := record[0]
		start, _ := strconv.ParseInt(record[1], 10, 64)
		end, _ := strconv.ParseInt(record[2], 10, 64)
		r := BEDBase{chr: chr, start: int(start), end: int(end)}
		set.record = append(set.record, &r)
	}
	return &set
}

func (bbs *BEDBaseSet) Write(path string) {
	file, err := os.Create(path)
	CheckError(err)
	defer file.Close()
	file.WriteString(bbs.String())
}

func (bbs *BEDBaseSet) String() string {
	s := strings.Builder{}
	for _, v := range bbs.record {
		s.WriteString(v.String() + "\n")
	}
	return s.String()
}

//func (bbs *BEDBaseSet) Sort() {
//	sort.Slice(bbs, func(i,j int) bool{return bbs.record[i].start > bbs.record[j].start})
//}

func (bbs *BEDBaseSet) Correct(bwfs *BEDWithFrequencySet, maxDistance float64, frequency int) {
	for _, bb := range bbs.record {
		var startArray []int
		var endArray []int
		//fmt.Println(bb)
		//time.Sleep(time.Second *4)
		for _, bwf := range bwfs.ele {
			//fmt.Println(bwf)
			switch {
			case bwf.score < frequency:
				continue
			case bwf.chr != bb.chr:
				continue
			case math.Abs(float64(bb.start-bwf.start)) < maxDistance:

				startArray = append(startArray, bwf.start)
			case math.Abs(float64(bb.end-bwf.start)) < maxDistance:

				endArray = append(endArray, bwf.start)
			}
		}
		// 计算最近的值
		bb.start = ProximityInt(bb.start, startArray)
		bb.end = ProximityInt(bb.end, endArray)

	}
}

// ProximityInt functions implements that find the proximity interger near point.
// If given spots have not a element, returns point value, immediately.
func ProximityInt(point int, spots []int) int {
	if len(spots) == 0 {
		return point
	}
	// 计算最近值所在的index
	type MapInt struct {
		value float64
		index int
	}
	var distances = make([]MapInt, len(spots))
	for i, v := range spots {
		distances[i] = MapInt{math.Abs(float64(point - v)), i}
	}
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].value < distances[j].value
	})
	return spots[distances[0].index]
}

type BEDsp struct { // Special BED format construction, including only 4 columns.
	BEDBase
	name string
}

func (b BEDsp) String() string {
	return fmt.Sprintf("%s\t%d\t%d\t%s", b.chr, b.start, b.end, b.name)
}

type FeatureSet struct {
	record []BEDsp
}

func (fs FeatureSet) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintln("# Special BED format\n# chr\tstart\tend\tname"))
	for _, r := range fs.record {
		s.WriteString(fmt.Sprintf("%s\n", r))
	}
	return s.String()
}

func (fs FeatureSet) Write(file *os.File) {
	defer file.Close()
	file.WriteString(fs.String())
}

// ReadFeatureSet function implements how to automatically read records from bed format files into FeatureSet struct.
// path refers to the path of bed format file.
func ReadFeatureSet(path string) *FeatureSet {
	var FS = FeatureSet{record: make([]BEDsp, 0)}
	f, err := os.Open(path)
	checkFile(err, path)
	defer f.Close()
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		next := scan.Text()
		if string(next[0]) != "#" {
			record := strings.Split(next, "\t")
			if len(record) == 4 {
				satrt, _ := strconv.ParseInt(record[1], 10, 64)
				end, _ := strconv.ParseInt(record[2], 10, 64)
				element := BEDsp{BEDBase{record[0], int(satrt), int(end)}, record[3]}
				FS.record = append(FS.record, element)
			} else {
				//log.Fatalln(path + "is not special BED format for BACtoPanda program")
				log.Println("new line\n" + next + "\n<<<<<<<<<<Look at, above!>>>>>>>>>>")
			}
		}
	}
	return &FS
}

type BEDWithFrequency struct {
	BEDsp
	// The score column in bed format specially refers to the frequency which each HindIII site occurs in sequencing reads.
	score int
}

func (bwf BEDWithFrequency) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%s\t%d\t%d\t%s\t%d", bwf.chr, bwf.start, bwf.end, bwf.name, bwf.score))
	return s.String()
}

type BEDWithFrequencySet struct {
	ele []BEDWithFrequency
}

func NewBEDWithFrequencySet() *BEDWithFrequencySet {
	return &BEDWithFrequencySet{ele: make([]BEDWithFrequency, 0)}
}

func ReadBEDWithFrequencySet(path string) *BEDWithFrequencySet {
	var bwfs = BEDWithFrequencySet{ele: make([]BEDWithFrequency, 0)}
	f, err := os.Open(path)
	checkFile(err, path)
	defer f.Close()
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		next := scan.Text()
		if string(next[0]) != "#" {
			record := strings.Split(next, "\t")
			//fmt.Printf("%v\t%d\n", record ,len(record))

			if len(record) == 5 {
				start, _ := strconv.ParseInt(record[1], 10, 64)
				end, _ := strconv.ParseInt(record[2], 10, 64)
				score, _ := strconv.ParseInt(record[4], 10, 64)
				element := BEDWithFrequency{BEDsp{BEDBase{record[0], int(start), int(end)},
					record[3]}, int(score)}
				bwfs.ele = append(bwfs.ele, element)
			} else {
				//log.Fatalln(path + "is not special BED format for BACtoPanda program")
				log.Println("new line\n" + next + "\n<<<<<<<<<<Look at, above!>>>>>>>>>>")
			}
		}
	}
	return &bwfs
}

func (bwfs BEDWithFrequencySet) String() string {
	s := strings.Builder{}
	for _, v := range bwfs.ele {
		s.WriteString(fmt.Sprintf("%s\n", v))
	}
	return s.String()
}

func (bwfs *BEDWithFrequencySet) Write(file *os.File) {
	file.WriteString(bwfs.String())
}

/*
	schedule：
		A: 把位点序列读入  pk
		B: 找到末端序列reads，并取末端序列 ok
		C: 通过末端序列找到位点，记录到字典中。map[string]BEDWithFrequency string = chr + start
*/

//var data = ReadFeatureSet("")

type BEDWithFrequencyMap struct {
	m map[string]*BEDWithFrequency
}

func NewBEDWithFrequencyMap() *BEDWithFrequencyMap {
	return &BEDWithFrequencyMap{
		m: make(map[string]*BEDWithFrequency, 0),
	}
}

func (nbwfm *BEDWithFrequencyMap) ToBEDWithFrequency() *BEDWithFrequencySet {
	bwfs := NewBEDWithFrequencySet()
	for _, v := range nbwfm.m {
		bwfs.ele = append(bwfs.ele, *v)
	}
	return bwfs
}

var mt sync.Mutex

func (nbwfm *BEDWithFrequencyMap) update(set *FeatureSet, seq string) {
	m := nbwfm.m
	for _, v := range set.record {
		/*
			ga := NewGlobalAligner(v.name[:len(seq)], seq)
			if ga.IsSimilar(4) {
		*/
		if equal(v.name[:len(seq)], seq) {
			mt.Lock() // For concurrency security
			key := fmt.Sprintf("%s%d", v.chr, v.start)
			if ele, ok := m[key]; ok {
				ele.score += 1
			} else {
				ele := BEDWithFrequency{
					BEDsp{
						BEDBase{
							v.chr, v.start, v.end,
						},
						"*"}, // 	用*代替序列
					1}
				m[key] = &ele
			}
			mt.Unlock()
		}
	}
}
