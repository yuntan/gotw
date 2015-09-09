package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/mitchellh/go-homedir"
	"github.com/yuntan/tw/go-tw"
)

var settingFile string

func main() {
	log.SetFlags(0)

	var err error
	if settingFile, err = homedir.Expand("~/.tw.json"); err != nil {
		log.Fatalf("error: failed to get home dir: %s\n", err)
	}

	if len(os.Args) == 1 {
		b, _ := ioutil.ReadAll(os.Stdin)

		setting := readSetting()
		if len(setting) == 0 {
			log.Fatalf("error: no account found\n")
		}

		for _, s := range setting {
			token := &oauth1.Token{
				Token:       s.AccessToken,
				TokenSecret: s.AccessTokenSecret,
			}
			resp, err := tw.Tweet(string(b), token)
			if err != nil || resp.StatusCode/100 != 2 {
				os.Exit(1)
			}
		}
	} else if len(os.Args) == 2 && os.Args[1] == "add" {
		addAccount()
	} else if len(os.Args) == 2 {
		setting := readSetting()
		if len(setting) == 0 {
			addAccount()
			setting = readSetting()
		}

		for _, s := range setting {
			token := &oauth1.Token{
				Token:       s.AccessToken,
				TokenSecret: s.AccessTokenSecret,
			}
			resp, err := tw.Tweet(os.Args[1], token)
			if err != nil {
				log.Fatalln(err)
			} else if resp.StatusCode/100 != 2 {
				log.Fatalln(resp.Status)
			} else {
				log.Println(resp.Status)
			}
		}
	} else if len(os.Args) == 4 && os.Args[1] == "fav" {
		// TODO
	}
}

func addAccount() {
	reqToken, url, err := tw.GetRequestToken()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Please open %s with your browser\n", url.String())
	fmt.Print("Enter PIN: ")
	var pin string
	fmt.Scanf("%s", &pin)

	token, err := tw.GetAccessToken(reqToken, pin)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := tw.VerifyCredentials(token)
	if err != nil {
		log.Fatalf("error: account verification failed: %s\n", err)
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	data := map[string]interface{}{}
	dec.Decode(&data)
	name := data["screen_name"].(string)

	appendSetting(name, AccountSetting{
		AccessToken:       token.Token,
		AccessTokenSecret: token.TokenSecret,
	})
	log.Printf("account %s saved\n", name)
}
