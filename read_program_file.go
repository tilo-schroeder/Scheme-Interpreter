package main

import (
	"bufio"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFromFile(filePath string) string {
	var sb strings.Builder
	data, err := os.Open(filePath)
	check(err)
	defer data.Close()

	s := bufio.NewScanner(data)
	for s.Scan() {
		read_line := s.Text()
		read_line = strings.TrimRight(read_line, "\r\n")
		sb.Write([]byte(read_line))
	}
	return sb.String()
}
