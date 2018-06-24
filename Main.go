package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
)

func main() {
	settings, err := readJSON()
	checkError(err)
	var node Node
	if settings.FirstInit {
		settings.FirstInit = false
		err = writeJSON(settings)
		checkError(err)
		node.firstInit()
	} else {
		node.init()
	}
}

func readJSON() (*JSONSettings, error) {
	data, err := ioutil.ReadFile("./Settings.json")
	if err != nil {
		return nil, err
	}
	var settings JSONSettings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		return nil, err
	}
	return &settings, err
}

func writeJSON(settings *JSONSettings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./Settings.json", data, 0644)
	return err
}

// JSONSettings is a data type for the json settings file
type JSONSettings struct {
	FirstInit  bool           `json:"FirstInit"`
	PrivateKey JSONPrivateKey `json:"PrivateKey"`
	Address    string         `json:"Address"`
}

// JSONPrivateKey is a data sub-type for the json settings file
type JSONPrivateKey struct {
	PublicKey JSONPublicKey `json:"PublicKey"`
	D         int64         `json:"D"`
}

// JSONPublicKey is a data sub-type for the json settings file
type JSONPublicKey struct {
	X int64 `json:"X"`
	Y int64 `json:"Y"`
}

// checkError calls panic() with the recieved error in case err != nil
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// getIPAddress returns the local ip address
func getIPAddress() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	str := localAddr.String()
	return str, nil
}
