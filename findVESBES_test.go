package main

import (
	"compress/gzip"
	"fmt"
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fastq"
	"github.com/biogo/biogo/seq/linear"
	"log"
	"os"
	"testing"
)

//var PROFILE =PROFILE("ACACACAC")

func TestFindOneVESBES(t *testing.T) {
	VESF = "CGCCATTCCTATGCGATGCAC"
	VESR = "ATAGACAAACTCAGATAC"
	destinationDir = "."

	f, err := os.Open("data_test/SP02-01_a-1r/SP02-01_a-1r_2.clean.fq.gz")
	if err != nil {
		log.Fatal("1112")
		log.Fatalln(err)
	}
	defer f.Close()

	uncompressFile, err := gzip.NewReader(f)
	CheckError(err)

	template := linear.NewSeq("", nil, alphabet.DNA)
	r := fastq.NewReader(uncompressFile, template)
	scan := seqio.NewScanner(r)
	for scan.Next() {
		fmt.Println(scan.Seq().Name())
	}
}

func TestFindTwoVESBES(t *testing.T) {
	VESF = "CGCCATTCCTATGCGATGCAC"
	VESR = "CGCCATAGCTATACGGCAC"
	destinationDir = "."
	FindTwoVESBES("data_test/SP02-01_a-1r")
}
