package manager

import (
	"os"
	"fmt"
)

func (dm Download) MergeFiles() error {
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
