package main

import "regexp"

func correctBED(samplePath string) {

	sampleIDRegexp := regexp.MustCompile(`SP[0-9]{2}-[0-9]{2}`)
	sampleID := sampleIDRegexp.FindStringSubmatch(samplePath)[0]

	bbs := ReadBEDBaseSet(string(destinationDir) + "/depthCalling/" + sampleID + "_BAC.bed")

	bwfs := ReadBEDWithFrequencySet(string(destinationDir) + "/VES-BES/" + sampleID + "_frequency.bed")

	bbs.Correct(bwfs, 5000, 5)
	bbs.Write(string(destinationDir) + "/depthCalling/" + sampleID + "_correct.bed")
}
