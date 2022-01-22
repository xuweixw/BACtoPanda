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
	for _, path := range sourceDir {
		pools, err := os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, pool := range pools {
			if pool.Name()[:2] == "SP" {
				samplePath := path + "/" + pool.Name()
				switch step {
				case ENTIRE:
					runBowtie2(samplePath)            //Bowtie2
					runSamtools(samplePath)           //Samtools
					calculatePotentialBAC(samplePath) //DepthCalling
					FindTwoVESBES(samplePath)
					correctBED(samplePath)
				case ALIGNMENT:
					runBowtie2(samplePath)
				case DEPTHCALLING:
					calculatePotentialBAC(samplePath)
				case EXTRACTION:
					FindTwoVESBES(samplePath)
				case CORRECT:
					correctBED(samplePath)
				}
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
