package manager

import (
	"fmt"
	"net/http"
	"io"
	"os"
)

func (dm Download) DownloadChunk(idx int, sec [2]int) error {
	r, err := dm.GetNewRequest("GET")
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