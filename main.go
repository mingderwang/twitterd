package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"os"
	"time"
)

func main() {
	var CONSUMER_KEY = os.Getenv("CONSUMER_KEY")
	var CONSUMER_SECRET = os.Getenv("CONSUMER_SECRET")
	var ACCESS_TOKEN = os.Getenv("ACCESS_TOKEN")
	var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api := anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)
	fmt.Println(*api.Credentials)
	api.EnableThrottling(1*time.Second, 5)

	searchResult, _ := api.DirectMessages(nil)
	for _, tweet := range searchResult.Statuses {
		fmt.Print(tweet.Text)
	}
}
