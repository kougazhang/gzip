package gzip

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(dir, "0.gz")
	write, err := NewWriter(old)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range strings.Split("1,2,3,4,5,6,7,8,9,10", ",") {
		_, err = write.Write(i)
		if err != nil {
			t.Fatal(err)
		}
	}
	if err = write.Close(); err != nil {
		t.Fatal(err)
	}

	reader, err := NewReader(old)
	if err != nil {
		t.Fatal(err)
	}

	for {
		line, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				_ = os.Remove(old)
				return
			}
		}
		fmt.Println("line: ", line)
	}
}

func TestNewWriter(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(dir, "0.gz")
	write, err := NewWriter(old)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range strings.Split("1,2,3,4,5,6,7,8,9,10", ",") {
		_, err = write.Write(i)
		if err != nil {
			t.Fatal(err)
		}
	}
	if err = write.Close(); err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(old)
}

func TestLimitedW_Renew(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	files := make([]string, 0, 4)
	old := filepath.Join(dir, "0.gz")
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
			t.Log("fullFile: ", fullFile)
			files = append(files, fullFile)
		}
	}
	if err = limit.Close(); err != nil {
		t.Fatal(err)
	}
	files = append(files, limit.Path)
	for _, file := range files {
		t.Log("delete: ", file)
		if err := os.Remove(file); err != nil {
			t.Fatal(err)
		}
	}
}
