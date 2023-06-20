package gzip

import (
	"bufio"
	"compress/gzip"
	"os"
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
	F  *os.File
	Gf *gzip.Writer
	Wf *bufio.Writer
}

func (w W) Close() error {
	if err := w.Wf.Flush(); err != nil {
		return err
	}
	if err := w.Gf.Flush(); err != nil {
		return err
	}
	if err := w.F.Close(); err != nil {
		return err
	}
	return w.Gf.Close()
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
		F:  f,
		Gf: gf,
		Wf: bufio.NewWriter(gf),
	}, nil
}
