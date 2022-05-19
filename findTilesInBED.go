package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Tile struct {
	chr             string
	start, end, len int64
	count           int64 //depth accelaration
	depth           float32
	next            *Tile
}

type TileLinkedlist struct {
	head *Tile
	tail *Tile
	len  int
}

func (tll TileLinkedlist) String() string {
	s := fmt.Sprintf("%v ... -> ... %v len=%d\n", tll.head, tll.tail, tll.tail.end-tll.head.start)
	return s
}

func (tll TileLinkedlist) ToBed() string {
	// chr	start	end
	s := fmt.Sprintf("%s\t%d\t%d\n", tll.head.chr, tll.head.start, tll.tail.end)
	return s
}

func (tll *TileLinkedlist) Add(t *Tile) {
	if tll.head == nil {
		tll.head = t
		tll.tail = t
		tll.len++
	} else {
		tll.tail.next = t
		tll.tail = t
		tll.len++
	}
}
func (tll TileLinkedlist) Len() int { //重新计算链表长度
	tll.len = 0
	for tll.head != nil {
		tll.len++
		tll.head = tll.head.next
	}
	return tll.len
}

func (tll TileLinkedlist) potentialBAC() bool {
	if tll.tail.end-tll.head.start > 10000 { //考虑Gap的存在，BAC最小长度设置为10Kb.
		return true
	}
	return false
}

func (tll TileLinkedlist) Splite(max int64) []TileLinkedlist {
	var tlls []TileLinkedlist
	current := tll
	current.len = 0
	for tll.head.next != nil {
		previous := tll.head
		tll.head = tll.head.next
		if tll.head.start-previous.end > max { //满足限制条件，断开链接
			current.tail = previous
			current.tail.next = nil
			current.len = current.Len()
			tlls = append(tlls, current)
			current = tll
		}
	}
	return tlls
}

type position struct {
	chr      string
	position int64
	depth    int64
}

func NewTiles(file *os.File) *TileLinkedlist {
	var p = new(position)
	var tile = new(Tile)
	var TileLL TileLinkedlist
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		newline := strings.Split(reader.Text(), "\t")
		p.chr = newline[0]
		p.position, _ = strconv.ParseInt(newline[1], 10, 64)
		p.depth, _ = strconv.ParseInt(newline[2], 10, 64)
		if tile.start == 0 { //判断首次传入数据
			tile.chr = p.chr
			tile.start = p.position
			tile.end = p.position
			tile.count = p.depth
		} else if tile.end+1 == p.position {
			tile.end = p.position
			tile.count += p.depth
		} else if tile.end+1 != p.position {
			TileLL.Add(tile)
			tile = new(Tile)
			tile.chr = p.chr
			tile.start = p.position
			tile.end = p.position
			tile.count = p.depth
		} else {
			log.Printf("have a issue for %v", p)
		}
	}
	return &TileLL
}

func calculatePotentialBAC(samplePath string) {
	sampeIDRegexp := regexp.MustCompile(`SP[0-9]{2}-[0-9]{2}`)
	sampleID := sampeIDRegexp.FindStringSubmatch(samplePath)[0]
	// `samtools depth -o mapping_result/SP01-01_depth.tab mapping_result/SP01-01_sorted.bam
	//generate depth.tab
	cmdDepth := exec.Command("samtools", "depth",
		"-o", string(destinationDir)+"/depthCalling/"+sampleID+"_depth.bed",
		string(destinationDir)+"/alignment/"+sampleID+"_sorted.bam")
	stdoutStderrDepth, err := cmdDepth.CombinedOutput()
	if err != nil {
		log.Fatalln("Running Samtools Depth generates an Error: ", err)
	} else {
		log.Println(sampleID, "ouput:\n", stdoutStderrDepth)
	}

	file, err := os.Open(string(destinationDir) + "/depthCalling/" + sampleID + "_depth.bed")
	defer file.Close()
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	finalOutput, err := os.Create(string(destinationDir) + "/depthCalling/" + sampleID + "_BAC.bed")
	defer finalOutput.Close()
	if err != nil {
		log.Fatalln("Creating Bed file counters a Error: ", err)
	}
	output := NewTiles(file)

	for _, item := range output.Splite(1000) {
		if item.potentialBAC() {
			finalOutput.WriteString(fmt.Sprintf("%v", item.ToBed()))
		}
	}
}
