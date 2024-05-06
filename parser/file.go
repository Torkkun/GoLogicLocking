package parser

import (
	"bufio"
	"os"
)

func NewReader(filename string) *os.File {
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return fp
}

func NewScanner(fp *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(fp)
	return scanner
}
