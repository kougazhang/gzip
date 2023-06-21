package gzip

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLimitedW_Renew(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	files := make([]string, 0, 4)
	old := filepath.Join(dir, "0.gz")
	files = append(files, old)
	limit, err := NewLimitedWriter(old, 3)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range strings.Split("1,2,3,4,5,6,7,8,9,10", ",") {
		_, err = limit.Write(i)
		if err != nil {
			t.Fatal(err)
		}
		if ok, _ := limit.IsFull(); ok {
			fullFile, err := limit.Renew()
			if err != nil {
				t.Fatal(err)
			}
			t.Log(fullFile)
			files = append(files, fullFile)
		}
	}
	limit.Close()
	files = append(files, limit.Path)
	for _, file := range files {
		os.Remove(file)
	}
}
