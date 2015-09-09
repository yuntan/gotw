package tw

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
)

const TWITTER_API = "https://api.twitter.com/1.1/"

var conf *oauth1.Config

func init() {
	conf = &oauth1.Config{
		ConsumerKey:    CONSUMER_KEY,
		ConsumerSecret: CONSUMER_SECRET,
		CallbackURL:    "oob",
	}
}

func GetRequestToken() (*oauth1.RequestToken, *url.URL, error) {
	conf.Endpoint = twitter.AuthorizeEndpoint
	reqToken, err := conf.GetRequestToken()
	if err != nil {
		return nil, nil, err
	}
	url, err := conf.AuthorizationURL(reqToken)
	if err != nil {
		return nil, nil, err
	}
	return reqToken, url, nil
}

func GetAccessToken(reqToken *oauth1.RequestToken, pin string) (*oauth1.Token, error) {
	return conf.GetAccessToken(reqToken, pin)
}

// https://dev.twitter.com/rest/reference/get/account/verify_credentials
func VerifyCredentials(token *oauth1.Token) (*http.Response, error) {
	conf.Endpoint = twitter.AuthenticateEndpoint
	client := oauth1.NewClient(conf, token)
	return client.Get(TWITTER_API + "account/verify_credentials.json")
}

// https://dev.twitter.com/rest/reference/post/statuses/update
func Tweet(tweet string, token *oauth1.Token) (*http.Response, error) {
	conf.Endpoint = twitter.AuthenticateEndpoint
	client := oauth1.NewClient(conf, token)
	return client.PostForm(TWITTER_API+"statuses/update.json",
		url.Values{"status": []string{tweet}})
}

// https://dev.twitter.com/rest/reference/get/statuses/mentions_timeline
func GetMentions(token *oauth1.Token) (*http.Response, error) {
	conf.Endpoint = twitter.AuthenticateEndpoint
	client := oauth1.NewClient(conf, token)
	return client.PostForm(TWITTER_API+"statuses/mentions_timeline.json", nil)
}

// https://dev.twitter.com/rest/reference/post/statuses/destroy/%3Aid
func RemoveTweet(id string, token *oauth1.Token) (*http.Response, error) {
	conf.Endpoint = twitter.AuthenticateEndpoint
	client := oauth1.NewClient(conf, token)
	return client.PostForm(TWITTER_API+fmt.Sprintf("statuses/destroy/%s.json", id), nil)
}
