package main

import (
	"fmt"
	"log"
	"time"
	"github.com/sepehrmohseni/go-cli-dm/manager"
)

func main() {
	var fileURL string
	var fileTarget string
	fmt.Println("Welcome to GO CLI Download Manager!")
	fmt.Println("What do you want to download?? (Input URL): ")
	fmt.Scanln(&fileURL)
	fmt.Println("Please enter absolute path for output file: ")
	fmt.Scanln(&fileTarget)
	startTime := time.Now()
	dm := manager.Download{
		Url: fileURL,
		TargetPath: fileTarget,
		TotalSections: 10,
	}
	if err := dm.Do(); err != nil {
		log.Printf("Some error occured: %s", err)
	}
	endTime := time.Now()
	completedTime := endTime.Sub(startTime).Seconds()
	fmt.Printf("Download completed in %v seconds. \n", completedTime)
}