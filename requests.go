package main

import (
	"bytes"
	"strconv"
)

// FormatSCM formats the th received hash and index to bytes
func FormatSCM(index int, hash string) []byte {
	str := strconv.Itoa(index) + "\000" + hash
	return []byte(str)
}

// UnformatSCM unformats the received bytes back to hash and index
func UnformatSCM(data []byte) (int, string, error) {
	splat := bytes.Split(data, []byte("\000"))
	index, err := strconv.Atoi(string(splat[0]))
	if err != nil {
		return 0, "", err
	}
	return index, string(splat[1]), nil
}

// FormatFT formats n to bytes
func FormatFT(blocks int) []byte {
	return []byte(strconv.Itoa(blocks))
}

// UnformatFT formats bytes to n
func UnformatFT(data []byte) (int, error) {
	blocks, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return blocks, nil
}

// FormatIS formats index to bytes
func FormatIS(index int) []byte {
	return []byte(strconv.Itoa(index))
}

// UnformatIS formats bytes to index
func UnformatIS(data []byte) (int, error) {
	index, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return index, nil
}

// FormatBP formats blocks array to byte array
func FormatBP(blocks []*Block) []byte {
	var data []byte
	for _, v := range blocks {
		bytes, err := v.ToBytes()
		if err != nil {
			continue
		}
		data = append(append(data, bytes...), []byte("|\000")...)
	}
	return data
}

// UnformatBP formats byte array to blocks array
func UnformatBP(data []byte) ([]*Block, error) {
	splat := bytes.Split(data, []byte("|\000"))
	blocks := make([]*Block, 0)
	for _, v := range splat {
		b := &Block{}
		b, err := ToBlock(v)
		if err != nil {
			continue
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

// FormatPA formats address (string) array to byte array
func FormatPA(addresses []string) []byte {
	var data []byte
	for _, addr := range addresses {
		data = append(append(data, []byte(addr)...), []byte("/\000")...)
	}
	return data
}

// UnformatPA formats byte array to address (string) array
func UnformatPA(data []byte) []string {
	splat := bytes.Split(data, []byte("/\000"))
	addrs := []string{}
	for _, v := range splat {
		addrs = append(addrs, string(v))
	}
	return addrs
}
