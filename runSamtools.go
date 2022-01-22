package main

import (
	"log"
	"os/exec"
	"regexp"
)

func runSamtools(samplePath string) {
	// parse filePrefix and sampleID
	//filePrefixRegexp := regexp.MustCompile(`SP01-[0-9]{2}_.*-1.{1}`)
	//filePrefix := filePrefixRegexp.FindStringSubmatch(samplePath)[0]
	sampeIDRegexp := regexp.MustCompile(`SP[0-9]{2}-[0-9]{2}`)
	sampleID := sampeIDRegexp.FindStringSubmatch(samplePath)[0]
	//launch Sample filte and sort
	//samtools view -f 2 -b -@ 100 SP01-01.sam | samtools sort -@ 100 > SP01-01.bam
	cmdView := exec.Command("samtools", "view",
		"-f", "2",
		"-b",
		"-@", string(threads),
		"-o", string(destinationDir)+"/alignment/"+sampleID+".bam",
		string(destinationDir)+"/alignment/"+sampleID+".sam")
	cmdSort := exec.Command("samtools", "sort",
		"-@", string(threads),
		"-o", string(destinationDir)+"/alignment/"+sampleID+"_sorted.bam",
		string(destinationDir)+"/alignment/"+sampleID+".bam")
	stdoutStderrView, err := cmdView.CombinedOutput()
	if err != nil {
		log.Println(sampleID, " View have an error: ", err)
	}
	log.Printf("%s View output: \n %s\n", sampleID, stdoutStderrView)
	stdoutStderrSort, err := cmdSort.CombinedOutput()
	if err != nil {
		log.Println(sampleID, " Sort have an error: ", err)
	}
	log.Printf("%s Sort output: \n%s\n", sampleID, stdoutStderrSort)
}
