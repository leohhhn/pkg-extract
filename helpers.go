package main

import (
	"encoding/json"
	"strings"
)

type TX struct {
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

func (tx TX) MarshalMetadata() ([]byte, error) {
	data := map[string]interface{}{
		"creator": tx.Msg[0].Creator,
		"deposit": tx.Msg[0].Deposit,
		// add what is needed
	}
	return json.MarshalIndent(data, "", "	")
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
