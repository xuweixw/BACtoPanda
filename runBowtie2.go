package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

//传入样本目录名称，完成自动分析，结果保存在mapping_result目录下
func runBowtie2(samplePath string) {
	// parse filePrefix and sampleID

	// filePrefixRegexp := regexp.MustCompile(`SP[0-9]{2}-[0-9]{2}_.*-1.{1}`)
	filePrefixRegexp := regexp.MustCompile(`SP[0-9]{2}-[0-9]{2}.*`)
	filePrefix := filePrefixRegexp.FindStringSubmatch(samplePath)[0]
	sampeIDRegexp := regexp.MustCompile(`SP[0-9]{2}-[0-9]{2}`)
	sampleID := sampeIDRegexp.FindStringSubmatch(samplePath)[0]
	//launch Bowtie2 mapping
	cmd := exec.Command("/home/Xuwei/miniconda3/bin/bowtie2",
		"--threads", threads.String(), "-x", string(genome),
		"-1", samplePath+"/"+filePrefix+"_1.clean.fq.gz", "-2", samplePath+"/"+filePrefix+"_2.clean.fq.gz", "-S", string(destinationDir)+"/alignment/"+sampleID+".sam") //路径需要写全 参数字符串前后不能有空格
	stdoutStderr, err := cmd.CombinedOutput() //结果输出到stderror而不是stdout
	if err != nil {
		log.Printf("Have a error in %s: %s\n", cmd.String(), err)
	}
	// extract the total of paired-end reads and the percent of alignment from standard output
	informationRegexp := regexp.MustCompile(`([0-9]*) \(.*%\) were paired[.|\n]*([0-9]{2}\.[0-9]{2})% overall alignment rate`)
	information := informationRegexp.FindStringSubmatch(fmt.Sprintf("%s", stdoutStderr))
	if information == nil {
		log.Println("parsing reads information don't succeed")
		log.Printf("%s\n", stdoutStderr)
	} else {
		statsTab(fmt.Sprintf("%s\t%s\t%s\n", sampleID, information[1], information[2]))
	}
	verboseLog(fmt.Sprintf("%s result:\n%s\n", sampleID, stdoutStderr))
}

func verboseLog(record string) {
	file := openLog(string(destinationDir) + "/alignment/verbose.log")
	defer file.Close()
	file.WriteString(record)
}

func statsTab(record string) {
	file := openLog(string(destinationDir) + "/alignment/stats.tab")
	defer file.Close()
	file.WriteString(record)
}
func openLog(s string) *os.File {
	var file *os.File
	var err error
	if file, err = os.Open(s); os.IsNotExist(err) {
		file, err = os.Create(s)
		if err != nil {
			log.Fatalf("creating %s file error:%s", s, err)
		}
	}
	return file
}
