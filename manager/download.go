package manager

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

	r, err := dm.getNewRequest("HEAD")
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

	var sections = make([][2]int, dm.TotalSections)
	eachSize := size / dm.TotalSections
	fmt.Printf("each chunk size is %v bytes\n", eachSize)

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

	var wg sync.WaitGroup
	for i, s := range sections {
		wg.Add(1)
		go func(i int, s [2]int) {
			defer wg.Done()
			err := dm.downloadChunk(i, s)
			if err != nil {
				panic(err)
			}
		}(i, s)
	}

	wg.Wait()
	err = dm.mergeFiles()
	if err != nil {
		return err
	}

	return nil
}

func (dm Download) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		dm.Url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", "Sep DM V1")
	return r, nil
}

func (dm Download) downloadChunk(idx int, sec [2]int) error {
	r, err := dm.getNewRequest("GET")
	if err != nil {
		return err
	}

	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", sec[0], sec[1]))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("error while processing %v", resp.StatusCode)
	}

	fmt.Printf("Downloaded %v bytes for Section %v\n", resp.ContentLength, idx)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("section-%v.tmp", idx), b, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

func (dm Download) mergeFiles() error {
	fileX, err := os.OpenFile(dm.TargetPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer fileX.Close()

	for i := 0; i < dm.TotalSections; i++ {
		fileTmp := fmt.Sprintf("section-%v.tmp", i)
		b, err := os.ReadFile(fileTmp)

		if err != nil {
			return err
		}

		bw, err := fileX.Write(b)
		if err != nil {
			return err
		}
		err = os.Remove(fileTmp)
		if err != nil {
			return err
		}
		fmt.Printf("Merged %v bytes from Section %v\n", bw, i)
	}

	return nil
}
