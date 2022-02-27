package main

import (
	"log"
	"os"
)

func main() {
	// flag parsing
	FlagParse()

	// check output directories
	// ---- destinationDir
	// ----|---- alignment
	checkSubdirectory(string(destinationDir), "alignment")
	// ----|---- depthCalling
	checkSubdirectory(string(destinationDir), "depthCalling")
	// ----|---- assignment
	checkSubdirectory(string(destinationDir), "assignment")
	// ----|---- stat
	checkSubdirectory(string(destinationDir), "stat")
	// ----/---- VES-BES
	checkSubdirectory(string(destinationDir), "VES-BES")

	// Alignment and DepthCalling
	for _, subDir := range sourceDir {
		pool, err := os.Stat(subDir)
		if err != nil {
			log.Fatal(err)
		}

		if pool.Name()[:2] == "SP" {
			switch step {
			case ENTIRE:
				runBowtie2(subDir)            //Bowtie2
				runSamtools(subDir)           //Samtools
				calculatePotentialBAC(subDir) //DepthCalling
				FindTwoVESBES(subDir)
				correctBED(subDir)
			case ALIGNMENT:
				runBowtie2(subDir)
			case DEPTHCALLING:
				calculatePotentialBAC(subDir)
			case EXTRACTION:
				FindTwoVESBES(subDir)
			case CORRECT:
				correctBED(subDir)
			}
		}
	}

	// Assignment
	switch step {
	case ENTIRE:
		Assignment()
	case ASSIGNMENT:
		Assignment()
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
