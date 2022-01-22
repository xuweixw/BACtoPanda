package main

import (
	"log"
	"os"
)

func checkSubdirectory(path, subDir string) {
	if _, err := os.Stat(path + "/" + subDir); os.IsNotExist(err) {
		os.Mkdir(path+"/"+subDir, 0700)
		log.Printf("A new subdirectory \"%s\" is created in %s\n", subDir, path)
	} else {
		log.Printf("The \"%s\" subdirectory is exist -------> PASS\n", subDir)
	}
}

func checkFile(err error, path string) {
	if err != nil {
		log.Fatalf("<%s> %v\n", path, err)
	}
}

func checkError(err error) {
	if err != nil {
		log.Println("from Check error")
		log.Fatalln(err)
	}
}
