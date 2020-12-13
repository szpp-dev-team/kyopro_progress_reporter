package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/earlgray283/kyopro_progress_reporter/util"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var port string
var api *slack.Client
var members *[]Member
var channelID = "C01E67SR0FP"

func main() {
	var err error

	if err := godotenv.Load(".env"); err != nil {
		log.Println("not found .env file")
	}

	api = slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	if err := util.DownloadFile("members.json"); err != nil {
		log.Fatal(err)
	}
	members, err = NewMemberFromJSON()
	if err != nil {
		log.Fatal(err)
	}
	if err := ConvertIdToName(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		if time.Now().Weekday() != time.Monday {
			return
		}
		if err := reportSubmissions(); err != nil {
			log.Println(err)
			if _, _, err := api.PostMessage(channelID, slack.MsgOptionText(fmt.Sprintf("エラーが起きたっピ！朗読するっピ！\n%s", err.Error()), false)); err != nil {
				log.Println(err)
			}
		}
	})

	http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			var res *slackevents.ChallengeResponse
			if err := json.Unmarshal(body, &res); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/plain")
			if _, err := w.Write([]byte(res.Challenge)); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		case slackevents.CallbackEvent:
			innerEvent := eventsAPIEvent.InnerEvent

			switch event := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				message := strings.Split(event.Text, " ")
				if len(message) < 2 {
					return
				}

				command := message[1]
				switch command {
				case "set":
					if _, _, err := api.PostMessage(event.Channel, slack.MsgOptionText("ここに json で設定を書き込む処理を入れるっぴ！", false)); err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusInternalServerError)
					}

				default:
					if _, _, err := api.PostMessage(event.Channel, slack.MsgOptionText("何を言っているのかがわかんないっぴ！ごめんっぴ！", false)); err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusInternalServerError)
					}
				}

				if _, _, err := api.PostMessage(event.Channel, slack.MsgOptionText("まだ週1報告以外の機能はできてないっぴ！ごめんなさいっぴ！", false)); err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	})

	log.Println("[INFO] Server listening...")
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal(err)
	}

}
