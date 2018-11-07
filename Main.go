package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"strings"
)

func main1() {
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

// getIPAddress returns the local ip address
func getIPAddress() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		if strings.Index(i.Name, "Wi-Fi") != 0 || strings.Index(i.Name, "eth0") != 0 {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {

			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP[0] == 0 {
					return v.IP, nil
				}
			case *net.IPAddr:
				return v.IP, nil
			}
		}
	}
	return nil, errors.New("IP not found")
}
