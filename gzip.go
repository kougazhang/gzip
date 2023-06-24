package gzip

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type R struct {
	F  *os.File
	Gf *gzip.Reader
	Rf *bufio.Reader
}

func (r R) Close() error {
	if err := r.F.Close(); err != nil {
		return err
	}
	return r.Gf.Close()
}

func (r R) ReadLine() (string, error) {
	return r.Rf.ReadString('\n')
}

func NewReader(path string) (*R, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	gf, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	return &R{
		F:  f,
		Gf: gf,
		Rf: bufio.NewReader(gf),
	}, nil
}

type W struct {
	Path string
	F    *os.File
	Gf   *gzip.Writer
	Wf   *bufio.Writer
}

func (w W) Close() error {
	if err := w.Wf.Flush(); err != nil {
		return err
	}
	if err := w.Gf.Flush(); err != nil {
		return err
	}
	if err := w.Gf.Close(); err != nil {
		return err
	}
	return w.F.Close()
}

func (w W) Write(s string) (int, error) {
	return w.Wf.WriteString(s)
}

func NewWriter(path string) (*W, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	gf := gzip.NewWriter(f)
	return &W{
		Path: path,
		F:    f,
		Gf:   gf,
		Wf:   bufio.NewWriter(gf),
	}, nil
}

// NewLimitedWriter
// path: the first file, the file pattern of subsequence files would like xx-seq-%d.gz,
// and the num would start from zero.
// maxLine: the max lines of one file
func NewLimitedWriter(path string, maxLine int64) (limit *LimitedW, err error) {
	limit = new(LimitedW)
	if limit.W, err = NewWriter(path); err != nil {
		return
	}
	limit.ext = filepath.Ext(path)
	limit.prefixPath = strings.Split(path, limit.ext)[0]
	limit.MaxLine = maxLine
	return
}

// LimitedW limited lines writer
type LimitedW struct {
	*W
	MaxLine, CurLine int64
	FileSeq          int
	prefixPath       string
	ext              string
}

// IsFull if the file is beyond of maxLine return true.
func (w *LimitedW) IsFull() (bool, error) {
	return w.CurLine >= w.MaxLine, nil
}

// Write a line with counter
func (w *LimitedW) Write(s string) (n int, err error) {
	n, err = w.W.Write(s)
	if err != nil {
		return
	}
	w.CurLine++
	return
}

// Renew close the old and create a new one
// the param path is the new file path
// the return value string is the old file path
func (w *LimitedW) Renew() (old string, err error) {
	old = w.Path
	if err = w.Close(); err != nil {
		return
	}

	newPath := w.prefixPath + fmt.Sprintf("-seq%d", w.FileSeq) + w.ext
	newW, err := NewWriter(newPath)
	if err != nil {
		return
	}
	w.FileSeq++

	w.W = newW
	w.CurLine = 0
	return
}
