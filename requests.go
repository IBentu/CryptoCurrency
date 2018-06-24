package main

import (
	"bytes"
	"strconv"
)

// formatSCM formats the th received hash and index to bytes
func formatSCM(index int, hash string) []byte {
	str := strconv.Itoa(index) + "\000" + hash
	return []byte(str)
}

// unformatSCM unformats the received bytes back to hash and index
func unformatSCM(data []byte) (int, string, error) {
	splat := bytes.Split(data, []byte("\000"))
	index, err := strconv.Atoi(string(splat[0]))
	if err != nil {
		return 0, "", err
	}
	return index, string(splat[1]), nil
}

// formatFT formats n to bytes
func formatFT(blocks int) []byte {
	return []byte(strconv.Itoa(blocks))
}

// unformatFT formats bytes to n
func unformatFT(data []byte) (int, error) {
	blocks, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return blocks, nil
}

// formatIS formats index to bytes
func formatIS(index int) []byte {
	return []byte(strconv.Itoa(index))
}

// unformatIS formats bytes to index
func unformatIS(data []byte) (int, error) {
	index, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return index, nil
}

func unformatNT() {
}
