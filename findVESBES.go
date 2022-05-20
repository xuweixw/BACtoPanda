package main

import (
	"fmt"
	"os"
	"regexp"
)

var (
	VESF = "TGCTGCAAGGCGATTAAGTTGGGTAACGCCAGGGTTTTCCCAGTCACGACGTTGTAAAAC" +
		"GACGGCCAGTGAATTGTAATACGACTCACTATAGGGCGAATTCGAGCTCGGTACCCGGGG" +
		"ATCCTCTAGAGTCGACCTGCAGGCATGC"
	VESR = "AGTTAGCTCACTCATTAGGCACCCCAGGCTTTACACTTTATGCTTCCGGCTCGTATGTTG" +
		"TGTGGAATTGTGAGCGGATAACAATTTCACACAGGAAACAGCTATGACCATGATTACGCC" +
		"AAGCTCTAATACGACTCACTATAGGGAGAC"
)

var (
	//Preparation for find HindIII site
	FS *FeatureSet
	// Preparation data structure
	BWFM *BEDWithFrequencyMap
)

// FindTwoVESBES function find forward and reverse VES-BES reads from Illumina Paired-end reads.
// The VES-BES reads will save in fastq format files.
func FindTwoVESBES(samplePath string) {
	filePrefixRegexp := regexp.MustCompile(`SP[0-9]{2}-[0-9]{2}.*`)
	filePrefix := filePrefixRegexp.FindStringSubmatch(samplePath)[0]
	sampleIDRegexp := regexp.MustCompile(`SP[0-9]{2}-[0-9]{2}`)
	sampleID := sampleIDRegexp.FindStringSubmatch(samplePath)[0]

	// Check file name suffix
	const (
		SUFFIX1 = ".clean.fq.gz"
		SUFFIX2 = ".fq.gz"
	)
	var suffix string = SUFFIX1
	if _, err := os.Stat(samplePath + "/" + filePrefix + "_1" + SUFFIX1); os.IsNotExist(err) {
		suffix = SUFFIX2
	}

	R1 := samplePath + "/" + filePrefix + "_1" + suffix
	R2 := samplePath + "/" + filePrefix + "_2" + suffix

	//Preparation for find HindIII site
	FS = ReadFeatureSet(PROFILE.String())
	// Preparation data structure
	BWFM = NewBEDWithFrequencyMap()

	// Forward
	outputF := string(destinationDir) + "/VES-BES/" + sampleID + "_forward.fasta"
	fmt.Println(R1)
	FindOneVESBES(R1, R2, VESF, outputF)

	// Reverse
	outputR := string(destinationDir) + "/VES-BES/" + sampleID + "_reverse.fasta"
	FindOneVESBES(R1, R2, VESR, outputR)

	// 保存结果
	bwf := BWFM.ToBEDWithFrequency()
	f, err := os.Create(string(destinationDir) + "/VES-BES/" + sampleID + "_frequency.bed")
	defer f.Close()
	checkError(err)
	bwf.Write(f)
}
