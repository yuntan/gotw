package main

import (
	"fmt"

	"github.com/dghubble/oauth1"
	"github.com/yuntan/tw/gotw"
)

func main() {
	if os.Exists("~/.tw.conf") {
		token = &oauth1.Token{
			AcsessToken:       "",
			AccessTokenSecret: "",
		}
	} else {
		reqToken, url, err := gotw.GetRequestToken()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("open %s\n", url.String())

		var pin string
		fmt.Scanf("%s", &pin)

		token, err := gotw.GetAccessToken(reqToken, pin)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	gotw.Tweet("test", token)
}
