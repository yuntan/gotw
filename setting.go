package main

import (
	"encoding/json"
	"log"
	"os"
)

type Setting map[string]AccountSetting

type AccountSetting struct {
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
	Header            string `json:"header"`
	Footer            string `json:"footer"`
}

func readSetting() Setting {
	r, err := os.OpenFile(settingFile, os.O_RDONLY, 0600)
	if err != nil {
		return Setting{}
	}
	defer r.Close()
	dec := json.NewDecoder(r)
	var data = Setting{}
	dec.Decode(&data)
	return data
}

func appendSetting(name string, s AccountSetting) {
	data := readSetting()
	if data == nil {
		data = Setting{}
	}

	data[name] = s
	w, _ := os.OpenFile(settingFile, os.O_CREATE|os.O_WRONLY, 0600)
	defer w.Close()
	enc := json.NewEncoder(w)
	// enc := json.NewEncoder(os.Stdout)
	err := enc.Encode(data)
	if err != nil {
		log.Fatalln(err)
	}
}
