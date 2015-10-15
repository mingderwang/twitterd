package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/k0kubun/twitter"
	"github.com/mingderwang/userstream"
	"github.com/parnurzeal/gorequest"
	"os"
	"strconv"
)

func main() {
	var CONSUMER_KEY = os.Getenv("CONSUMER_KEY")
	var CONSUMER_SECRET = os.Getenv("CONSUMER_SECRET")
	var ACCESS_TOKEN = os.Getenv("ACCESS_TOKEN")
	var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")
	client := &userstream.Client{
		ConsumerKey:       CONSUMER_KEY,
		ConsumerSecret:    CONSUMER_SECRET,
		AccessToken:       ACCESS_TOKEN,
		AccessTokenSecret: ACCESS_TOKEN_SECRET,
	}

	client.UserStream(func(event interface{}) {
		switch event.(type) {
		case *twitter.Tweet:
			tweet := event.(*twitter.Tweet)
			fmt.Printf("%s: %s\n", tweet.User.ScreenName, tweet.Text)
		case *userstream.Delete:
			tweetDelete := event.(*userstream.Delete)
			fmt.Printf("[delete] %d\n", tweetDelete.Id)
		case *userstream.Favorite:
			favorite := event.(*userstream.Favorite)
			fmt.Printf("[favorite] %s => %s : %s\n",
				favorite.Source.ScreenName, favorite.Target.ScreenName, favorite.TargetObject.Text)
		case *userstream.Unfavorite:
			unfavorite := event.(*userstream.Unfavorite)
			fmt.Printf("[unfavorite] %s => %s : %s\n",
				unfavorite.Source.ScreenName, unfavorite.Target.ScreenName, unfavorite.TargetObject.Text)
		case *userstream.Follow:
			follow := event.(*userstream.Follow)
			fmt.Printf("[follow] %s => %s\n", follow.Source.ScreenName, follow.Target.ScreenName)
		case *userstream.Unfollow:
			unfollow := event.(*userstream.Unfollow)
			fmt.Printf("[unfollow] %s => %s\n", unfollow.Source.ScreenName, unfollow.Target.ScreenName)
		case *userstream.ListMemberAdded:
			listMemberAdded := event.(*userstream.ListMemberAdded)
			fmt.Printf("[list_member_added] %s (%s)\n",
				listMemberAdded.TargetObject.FullName, listMemberAdded.TargetObject.Description)
		case *userstream.ListMemberRemoved:
			listMemberRemoved := event.(*userstream.ListMemberRemoved)
			fmt.Printf("[list_member_removed] %s (%s)\n",
				listMemberRemoved.TargetObject.FullName, listMemberRemoved.TargetObject.Description)
		case *userstream.Record:
			directMessage := event.(*userstream.Record)
			sendRequest(directMessage.DirectMessage.Sender.ScreenName, directMessage.DirectMessage.Text)
		}
	})
}

func stringify(jsonString string) string {
	str := strconv.Quote(jsonString)
	return str
}

func sendRequest(userName string, jsonSchema string) {
	fmt.Printf("%s : %s\n", userName, jsonSchema)
	request := gorequest.New()
	str := `{"domainName":"` + userName + `","typeName":"ggcoin","jsonSchema":` + stringify(jsonSchema) + `}`
	resp, _, _ := request.Post("http://log4security.com:8080/onion").
		Set("Content-Type", "application/json").
		Send(str).
		End()
	//spew.Dump(resp)
}
