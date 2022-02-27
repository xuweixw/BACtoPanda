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
)

//const border float64 = 3000

type Region struct {
	name       string
	chr        string
	start, end int64
}

func (r Region) String() string {
	return fmt.Sprintf("%s\t%d\t%d\t%s\n", r.chr, r.start, r.end, r.name)
}

// Equal method determines two Regions that are completely or extremely same.
// constant "border" defines the difference.
func (r Region) Equal(ele Region) (*Region, bool) {
	if r.chr == ele.chr {
		if math.Abs(float64(r.start-ele.start)) <= float64(margin) && math.Abs(float64(r.end-ele.end)) <= float64(margin) {
			var start, end int64
			if r.start > ele.start {
				start = ele.start
			} else {
				start = r.start
			}
			if r.end > ele.end {
				end = r.end
			} else {
				end = ele.end
			}
			r.start = start
			r.end = end
			return &r, true
		}
	}
	return nil, false
}

type RegionSet struct {
	regions []Region
}

// InitRegionSet function initialize RegionSet.
func InitRegionSet() *RegionSet {
	var rs RegionSet
	rs.regions = make([]Region, 0)
	return &rs
}
func (rs RegionSet) String() string {
	str := strings.Builder{}
	for _, item := range rs.regions {
		str.WriteString(fmt.Sprintf("%s", item))
	}
	return str.String()
}

func Intersection(a, b *RegionSet) *RegionSet {
	var set = RegionSet{regions: make([]Region, 0)}
	//fmt.Println("a::::",*a)
	//fmt.Println("b::::", *b)
	for _, av := range a.regions {
		for _, bv := range b.regions {
			if r, ok := av.Equal(bv); ok {
				set.regions = append(set.regions, *r)
			}
		}
	}
	return &set
}

func NewRegionSet(filePath string) *RegionSet {
	var set = InitRegionSet()
	file, err := os.Open(filePath)
	checkFile(err, filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := strings.Split(scanner.Text(), "\t")
		chr := record[0]
		start, _ := strconv.ParseInt(record[1], 10, 64)
		end, _ := strconv.ParseInt(record[2], 10, 64)
		r := Region{chr: chr, start: start, end: end}
		set.regions = append(set.regions, r)
	}
	return set
}

// location function assignment the BAC regions and return accurate BAC with name and the number of lose BAC.
func Location() (*RegionSet, int) {
	var rowMap = make(map[string]*RegionSet, 0) // Sixteen rows in a 384-well plate with letter form upper-case A to P.

	var colMap = make(map[string]*RegionSet, 0) //
	var poolMap = make(map[int]*RegionSet, 0)
	var results = InitRegionSet()
	var num int
	for i := 1; i <= 25; i++ {
		poolMap[i] = NewRegionSet(fmt.Sprintf(string(destinationDir)+"/depthCalling/%s-%02d_correct.bed", UNIT, i))
	}
	for k, v := range ROWID {
		rowMap[k] = Intersection(poolMap[v[0]], poolMap[v[1]])
	}
	for k, v := range COLID {
		colMap[k] = Intersection(poolMap[v[0]], poolMap[v[1]])
	}

	for key, p := range UNITS[string(UNIT)] {
		for r, _ := range ROWID {
			for c, _ := range COLID {
				BACID := fmt.Sprintf("%s%s%s", p, r, c)

				plate := Intersection(rowMap[r], colMap[c])
				well := Intersection(poolMap[key+1], plate)
				if len(well.regions) == 0 { // lose stat
					num++
				} else {
					for i, _ := range well.regions {
						well.regions[i].name = BACID
						results.regions = append(results.regions, well.regions[i])
					}
				}
			}
		}
	}
	sort.Slice(results.regions, func(i, j int) bool { return results.regions[i].chr < results.regions[j].chr })
	return results, num
}

func Assignment() {
	// Assignment
	r, num := Location()
	file, err := os.Create(string(destinationDir) + "/assignment/" + string(UNIT) + "_BAC.bed")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	if n, err := file.WriteString(fmt.Sprintf("%s", r)); err != nil {
		log.Printf("encounter an issue in writing : %s\n", err)
	} else {
		log.Printf("There are %d BAC parsed succussfully.\n\t\tThere are %d chars written", 7*384-num, n)
	}
}
