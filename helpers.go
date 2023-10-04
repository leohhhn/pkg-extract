package main

import (
	"bytes"
	"io"
	"strings"
)

type LogLine struct {
	Msg []struct {
		Type    string `json:"@type"`
		Creator string `json:"creator"`
		Package struct {
			Name  string `json:"Name"`
			Path  string `json:"Path"`
			Files []struct {
				Name string `json:"Name"`
				Body string `json:"Body"`
			} `json:"Files"`
		} `json:"package"`
		Deposit string `json:"deposit"`
	} `json:"msg"`
}

func findSubstringPositions(input, substring string) []int {
	var positions []int
	start := 0

	for {
		index := strings.Index(input[start:], substring)
		if index == -1 {
			break
		}
		position := start + index
		positions = append(positions, position)
		start = position + len(substring)
	}

	return positions
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
