package main

import (
	"log"
	"os"
	"strings"

	"github.com/jmcvetta/randutil"
	"github.com/nlopes/slack"
)

const (
	API_KEY    = "<your_api-key>"           // your api_key
	MENTION    = "<your_slack_bot_mention>" // your slack_bot mention
	NOT_TARGET = "not_target"
	HUNGRY     = "おなかすいた"
)

type Resutaurant struct {
	Name   string
	Weight int
}

func getRestaurants() string {
	choices := make([]randutil.Choice, 0, 2)
	choices = append(choices, randutil.Choice{10, "みんな"})
	choices = append(choices, randutil.Choice{5, "いっちゃん"})
	choices = append(choices, randutil.Choice{5, "ユキちゃん"})
	choices = append(choices, randutil.Choice{4, "あっちゃん"})
	choices = append(choices, randutil.Choice{3, "中村屋"})
	choices = append(choices, randutil.Choice{3, "りんご"})
	choices = append(choices, randutil.Choice{2, "今日はスーパーで弁当買いますか"})
	choices = append(choices, randutil.Choice{1, "なし"})

	result, _ := randutil.WeightedChoice(choices)
	return result.Item.(string)
}

func run(api *slack.Client) int {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				res := handleMessageEvent(ev)
				if res != NOT_TARGET {
					rtm.SendMessage(rtm.NewOutgoingMessage(res, ev.Channel))
				}
			}
		}
	}
}

func handleMessageEvent(ev *slack.MessageEvent) string {
	log.Printf("[INFO] User name: " + ev.Msg.Text)
	// bot宛てではない場合無視
	// if strings.Index(ev.Msg.Text, MENTION) < 0 {
	//   return "", NOT_TARGET
	// }
	// おなかすいた
	if strings.Index(ev.Msg.Text, HUNGRY) >= 0 {
		return getRestaurants()
	}
	return NOT_TARGET
}

func main() {
	api := slack.New(API_KEY)
	os.Exit(run(api))
}
