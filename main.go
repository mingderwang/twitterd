package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/davecgh/go-spew/spew"
	"github.com/k0kubun/twitter"
	"github.com/mingderwang/userstream"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	EndpointURL string `json:"endpointURL"`
}

type Onion struct {
	Ginger_Created int32 `json:"ginger_created"`
	Ginger_Id      int32 `json:"ginger_id" gorm:"primary_key"`

	DomainName string `json:"domainName"`
	TypeName   string `json:"typeName"`
	JsonSchema string `json:"jsonSchema"`
}

const (
	errorString = "OOPS! format error, correct example:\n {\"user\":{\"name\":\"John\", \"age\":32}}"
)

var (
	baseURL = "http://log4security.com:8080/onion"
)

var CONSUMER_KEY = os.Getenv("CONSUMER_KEY")
var CONSUMER_SECRET = os.Getenv("CONSUMER_SECRET")
var ACCESS_TOKEN = os.Getenv("ACCESS_TOKEN")
var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")

func main() {
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
			//spew.Dump(directMessage.DirectMessage.Text)
			str := directMessage.DirectMessage.Text
			fmt.Println("-------")
			spew.Dump(str)
			fmt.Println("-------")
			if len(str) != 0 {
				b := []byte(str)
				var f map[string]interface{}
				err := json.Unmarshal(b, &f)
				if err != nil && str != errorString {
					fmt.Println("-2------")
					log.Println(err)
					callBackUser(directMessage.DirectMessage.Sender.ID, "OOPS! format error, correct example:\n {\"user\":{\"name\":\"John\", \"age\":32}}")
				} else if len(f) == 1 {
					s2 := str
					s2 = strings.TrimSuffix(s2, "}")
					fmt.Println("s2:", s2)
					s2 = strings.TrimPrefix(s2, "{")
					fmt.Println("s2:", s2)
					sendRequest(directMessage.DirectMessage.Sender.ScreenName, directMessage.DirectMessage.Sender.ID, s2)
				}
			}
		}
	})
}

func stringify(data string) (tag string, schema string) {
	jsonString := strings.SplitN(data, ":", 2)
	//	spew.Dump(jsonString)
	if len(jsonString) == 2 {
		str := strconv.Quote(jsonString[1])
		return jsonString[0], str
	} else {
		return "", ""
	}
}

func sendRequest(userName string, id int, jsonSchemaWithTag string) {
	request := gorequest.New()
	if tag, schema := stringify(jsonSchemaWithTag); tag == "" && schema == "" {
		fmt.Println("error")
	} else {
		str := `{"domainName":"` + userName + `","typeName":` + tag + `,"jsonSchema":` + schema + `}`

		fmt.Printf("%s", str)
		resp, _, err := request.Post(baseURL).
			Set("Content-Type", "application/json").
			Send(str).End()
		if err != nil {
			//spew.Dump(err)
			callBackUser(id, "JSON syntax error")
		} else {

			//	spew.Dump(body)
			//	spew.Dump(resp)
			target := Onion{}
			processResponser(resp, &target)
			spew.Dump(target.Ginger_Id)
			var s string = strconv.Itoa(int(target.Ginger_Id))
			sendRequestByIdForBuild(s, id)
		}
	}
}

func sendRequestByIdForBuild(idString string, id int) {
	target := Result{}
	url := fmt.Sprintf("%s/%s/build", baseURL, idString)
	fmt.Println(url)
	request := gorequest.New()
	resp, _, _ := request.Get(url).End(printStatus)
	//spew.Dump(resp.Body)
	json.NewDecoder(resp.Body).Decode(&target)
	//spew.Dump(target)
	time.Sleep(1000 * time.Millisecond)
	callBackUser(id, target.EndpointURL)
}

func callBackUser(id int, endpoint string) {
	spew.Dump(id)
	spew.Dump(endpoint)
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api := anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)
	_, err := api.PostDMToUserId(endpoint, int64(id))
	//spew.Dump(message)
	if err != nil {
		spew.Dump(err)
	}
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

func processResponser(response *http.Response, target *Onion) {
	json.NewDecoder(response.Body).Decode(&target)
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
