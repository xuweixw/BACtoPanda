package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// genenome sequence
type Genome string

func (g *Genome) String() string { return fmt.Sprintf("%s", *g) }
func (g *Genome) Set(val string) error {
	var err error
	if _, err := os.Stat(val); os.IsNotExist(err) {
		err = errors.New("Genome sequenc can't find out in given file path <" + val + ">. ")
	} else {
		*g = Genome(val)
	}
	return err
}

type EnzymeProfile string

func (ep *EnzymeProfile) String() string { return fmt.Sprintf("%s", *ep) }
func (ep *EnzymeProfile) Set(val string) error {
	var err error
	if _, err := os.Stat(val); os.IsNotExist(err) {
		err = errors.New("Enzyme profile file can't find out in given file path <" + val + ">. ")
	} else {
		*ep = EnzymeProfile(val)
	}
	return err
}

// -- Unit Value
// Create a flag with a Unit name, example, -unit=SP01
// implement flag.Value interface
type Unit string

func (p *Unit) String() string { return string(*p) }

// 检查指定的单元名是否在UNITS中登记
func (p *Unit) Set(val string) error {
	var err error
	if _, ok := UNITS[val]; ok{
		*p = Unit(val)
	} else {
		err = errors.New(val + " is not a valid unit name.")
	}
	return err
}

// -- Source Directory Value
// -- Check weather the disrectory is exist and all 25 pool sequence files in here
type SourceDirectory []string

func (s *SourceDirectory) String() string { return fmt.Sprintf("%s", *s) }

func (s *SourceDirectory) Set(val string) error {
	var err error
	numFile := 0
	for _, path := range strings.Split(val, ",") {
		if _, err := os.Stat(path); os.IsNotExist(err) { // check existence
			return errors.New("Given the source directory is not exist.")
		}
		if files, err := ioutil.ReadDir(path); err != nil { // check number of files in directory
			return err
		} else {
			fileNameRegexp := regexp.MustCompile(`SP[0-9]{2}-[0-9]{2}_.*-1.{1}`)
			for _, file := range files {
				name := file.Name()
				if fileNameRegexp.FindStringSubmatch(name) != nil {
					numFile++
				}
			}
		}
		*s = append(*s, path)
	}
	if numFile != 25 {
		log.Printf("Given the source directory <%s>, only %d sequence files not 25.", val, numFile)
		os.Exit(0)
	}
	return err
}

// -- Working Directory Value
type DestinationDirectory string

func (d *DestinationDirectory) String() string { return fmt.Sprintf("%s", *d) }
func (d *DestinationDirectory) Set(val string) error {
	var err error
	if _, err := os.Stat(val); os.IsNotExist(err) {
		if os.MkdirAll(val, 700) != nil {
			log.Printf("Create the given destination directory <%s> successfully.\n", val)
		} else {
			err = errors.New("Unfortunately, don't create the given destination directory <" + val + ">")
		}
	} else if os.IsExist(err) {
		log.Println("Given destination directory <" + val + "> is exist.")
	}
	return err
}

// -- Row Pool Layout
var RowLayoutDefault = "------- ||  SP00-12 |  SP00-13 |  SP00-14 |  SP00-15 |\n" +
	"------------------------------------------------------\n" +
	"SP00-08 || SP00-R-A | SP00-R-B | SP00-R-C | SP00-R-D |\n" +
	"SP00-09 || SP00-R-E | SP00-R-F | SP00-R-G | SP00-R-H |\n" +
	"SP00-10 || SP00-R-I | SP00-R-J | SP00-R-K | SP00-R-L |\n" +
	"SP00-11 || SP00-R-M | SP00-R-N | SP00-R-O | SP00-R-P |"

// -- Column Pools Layout
var ColumnLayoutDefault = "------- ||  SP00-21  |  SP00-22  |  SP00-23  |  SP00-24  |  SP00-25  |\n" +
	"----------------------------------------------------------------------\n" +
	"SP00-16 || SP00-C-01 | SP00-C-02 | SP00-C-03 | SP00-C-04 | SP00-C-05 |\n" +
	"SP00-17 || SP00-C-06 | SP00-C-07 | SP00-C-08 | SP00-C-09 | SP00-C-10 |\n" +
	"SP00-18 || SP00-C-11 | SP00-C-12 | SP00-C-13 | SP00-C-14 | SP00-C-15 |\n" +
	"SP00-19 || SP00-C-16 | SP00-C-17 | SP00-C-18 | SP00-C-19 | SP00-C-20 |\n" +
	"SP00-20 || SP00-C-21 | SP00-C-22 | SP00-C-23 | SP00-C-24 |    ---    |"

// -- Step Control
// The entire analysis workflow consists of Alignment, DepthCalling and Assignment
type StepControl int

const (
	ENTIRE = iota
	ALIGNMENT
	DEPTHCALLING
	ASSIGNMENT
	EXTRACTION
	CORRECT
)

func (s *StepControl) String() string { return fmt.Sprintf("%v", *s) }

func (s *StepControl) Set(val string) error {
	var err error
	switch val {
	case "alignment":
		*s = ALIGNMENT
	case "depthcalling":
		*s = DEPTHCALLING
	case "assignment":
		*s = ASSIGNMENT
	case "entire":
		*s = ENTIRE
	case "extraction":
		*s = EXTRACTION
	case "correct":
		*s = CORRECT
	default:
		err = errors.New("The argument of <step> flag must be " +
			"one of <entire>, <alignment>, <depthcalling>, <correct> and <assignment>")
	}
	return err
}

type Margin float64

func (m *Margin) String() string { return fmt.Sprintf("%.0f", *m) }
func (m *Margin) Set(val string) error {
	var err error
	if n, err := strconv.ParseFloat(val, 10); err != nil {
		err = errors.New("Error: the argument of margin flag should be a number\n")
		flag.Usage()
		os.Exit(0)
	} else if n < 0 {
		err = errors.New("Error: the argument of margin flag must be more than 0\n")
		flag.Usage()
		os.Exit(0)
	} else {
		*m = Margin(n)
	}
	return err
}

// -- Threads Value
type Threads uint8

func (t *Threads) String() string { return fmt.Sprintf("%d", *t) }

func (t *Threads) Set(val string) error {
	var err error
	n, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		return err
	}
	if int(n) <= runtime.NumCPU() {
		*t = Threads(n)
	} else {
		err = errors.New("Given the thread number is more than the maximum amount of CPU cores.")
	}
	return err
}

func setted(value flag.Value) bool {
	if len(value.String()) > 0 {
		return true
	} else {
		return false
	}
}

var (
	// genome sequence for bowtie2
	genome Genome

	PROFILE EnzymeProfile

	// UNIT variable stores each unit name.
	UNIT Unit

	// source directory
	sourceDir SourceDirectory

	//destination directory
	destinationDir DestinationDirectory

	// margin size for BAC intersection
	margin Margin = 200

	// step control
	step StepControl = ENTIRE

	// threads number
	threads Threads = 1

	// newest version
	version = "BACtoPanda v2.0.1"
)

func FlagParse() {
	// Desfault Value for destDir
	s, _ := os.Getwd()
	destinationDir = DestinationDirectory(s)
	// Get the arguments form command line
	flag.Var(&genome, "genome", "Specify the path of the genome sequence used in bowtie2 alignment")
	flag.Var(&PROFILE, "profile", "Specify the path of the enzyme profile file for accuracy of BAC terminals.")
	flag.Var(&sourceDir, "sourceDir", "Specify the source directory"+
		" in which all 25 files in a supperpool must be exist.")
	flag.Var(&destinationDir, "destDir", "Specify the destination directory for all outputs.")
	flag.Var(&UNIT, "unit", "Specify the palte numbers"+
		" with a comma-seperated string.\n\texample: 1,2,3,4,5,6,7")
	flag.String("rowLayout", "", RowLayoutDefault)       // show the default row pool layout, cann't change.
	flag.String("columnLayout", "", ColumnLayoutDefault) // show the default column pool layout, cann't change.
	flag.Var(&margin, "margin", "Specify the width of the margin using in BAC assignment")
	flag.Var(&step, "step", "Control analysis workflows."+
		" The argment must be one of <entire>, <alignment>, <depthcalling> \nand <assignment>, <extraction>, <correct>. (default <entire>)")
	flag.Var(&threads, "threads", "Specify the number of threads to use.")
	v := flag.Bool("v", false, "Print version information.")

	// parse
	if len(os.Args) > 1 {
		flag.Parse()
	} else {
		flag.Usage()
		os.Exit(0)
	}
	if *v {
		fmt.Println(version)
		os.Exit(0)
	}
	// Check mandatory arguments
	if !setted(&genome) {
		flag.Usage()
		log.Fatal("Argument of <genome> is missing")
	}
	if !setted(&PROFILE) {
		flag.Usage()
		log.Fatal("Argument of <profile> is missing")
	}
	if !setted(&sourceDir) {
		flag.Usage()
		log.Fatal("Argument of <sourceDir> is missing")
	}
	if !setted(&UNIT) {
		flag.Usage()
		log.Fatal("Argument of <plates> is missing")
	}
}
