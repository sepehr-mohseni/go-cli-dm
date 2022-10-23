package manager

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Download struct {
	Url           string
	TargetPath    string
	TotalSections int
}

func (dm Download) Do() error {
	fmt.Println("Please wait while validating the URL...")

	// head request
	r, err := dm.GetNewRequest("HEAD")
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	fmt.Printf("This is the response %v.\n\n", resp.StatusCode)

	if resp.StatusCode > 299 {
		return fmt.Errorf("error while processing %v", resp.StatusCode)
	}

	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	fmt.Printf("File size is %f MB\n\n", float64(size)/1000000.0)

	// create chunks of file
	var sections = make([][2]int, dm.TotalSections)
	eachSize := size / dm.TotalSections
	fmt.Printf("each chunk size is %v bytes\n", eachSize)

	// [[0 10][10 20]....[99 end-1]] - total 100 byte file
	for i := range sections {
		if i == 0 {
			sections[i][0] = 0
		} else {
			sections[i][0] = sections[i-1][1] + 1
		}

		if i == dm.TotalSections-1 {
			sections[i][1] = size - 1
		} else {
			sections[i][1] = sections[i][0] + eachSize
		}
	}
	log.Println(sections)

	// download each chunk concurrently
	var wg sync.WaitGroup
	for i, s := range sections {
		wg.Add(1)
		go func(i int, s [2]int) {
			defer wg.Done()
			err := dm.DownloadChunk(i, s)
			if err != nil {
				panic(err)
			}
		}(i, s)
	}

	wg.Wait()
	err = dm.MergeFiles()
	if err != nil {
		return err
	}

	return nil
}
