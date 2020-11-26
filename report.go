package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/earlgray283/kyopro_progress_reporter/atcoder"
	"github.com/slack-go/slack"
)

type UserACcount struct {
	member        Member
	ACcount       int
	UniqueACcount int
}

// var Homekotoba = []string{"静大の誇りっぴ〜！", "その調子っぴ！", "もうすこし頑張るっぴ！", "課題をするな競技プログラミングをしろ"}

func reportSubmissions() error {
	since := int64(1606057200) // 2020-11-16(Mon) 00:00:00
	until := since + 604800

	for {
		userACcount, err := generateRanking(since, until)
		if err != nil {
			log.Println(err)
		}

		sinceDay := time.Unix(since, 0).Format("01/02(Mon)")
		untilDay := time.Unix(until, 0).Format("01/02(Mon)")

		msg := fmt.Sprintf("%s ~ %s にかけての AC 数ランキングを発表するっぴ！\n\n", sinceDay, untilDay)
		for rank, user := range *userACcount {
			msg += fmt.Sprintf("\t%d位: <@%s> (AC count: %d(%d))\n", rank+1, user.member.SlackID, user.ACcount, user.UniqueACcount)
			msg += fmt.Sprintf("\t%s\n", func() string {
				if user.ACcount <= 30 {
					return "は？"
				} else if 31 <= user.ACcount && user.ACcount < 40 {
					return "その調子っぴ！"
				} else {
					return "静大の誇りっぴ〜！"
				}
			}())
			msg += "\n"
		}

		if _, _, err := api.PostMessage(channelID, slack.MsgOptionText(msg, false)); err != nil {
			log.Println(err)
		}

		time.Sleep(time.Hour)
		since += 3600
		until += 3600
	}
}

func generateRanking(since, until int64) (*[]UserACcount, error) {
	list := []UserACcount{}

	for _, user := range *members {
		ACcount, err := countAC(user.AtCoderID, since, until)
		if err != nil {
			return nil, err
		}
		UniqueACcount, err := countUniqueAC(user.AtCoderID, since, until)
		if err != nil {
			return nil, err
		}

		list = append(list, UserACcount{member: user, ACcount: ACcount, UniqueACcount: UniqueACcount})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ACcount > list[j].ACcount
	})

	return &list, nil
}

func countAC(name string, since, until int64) (int, error) {
	time.Sleep(2 * time.Second) // 2s ごとに api 叩く
	srs, err := atcoder.GetSubmissionResult(name)
	if err != nil {
		return 0, err
	}

	sort.Slice(*srs, func(i, j int) bool {
		return (*srs)[i].EpochSecond < (*srs)[j].EpochSecond
	})

	ok := int64(-1)
	ng := int64(len((*srs)))
	for math.Abs(float64(ok-ng)) > 1 {
		mid := (ok + ng) / 2
		if (*srs)[mid].EpochSecond == since {
			ok = mid
			break
		}
		if (*srs)[mid].EpochSecond > since {
			ng = mid
		} else {
			ok = mid
		}
	}
	ok += 1
	sum := 0
	for _, sr := range (*srs)[ok:] {
		if since <= sr.EpochSecond && sr.EpochSecond <= until {
			sum++
		}
	}

	return sum, nil
}

func countUniqueAC(name string, since, until int64) (int, error) {
	time.Sleep(2 * time.Second) // 2s ごとに api 叩く
	srs, err := atcoder.GetUniqueAC(name)
	if err != nil {
		return 0, err
	}

	sort.Slice(*srs, func(i, j int) bool {
		return (*srs)[i].EpochSecond < (*srs)[j].EpochSecond
	})

	ok := int64(-1)
	ng := int64(len((*srs)))
	for math.Abs(float64(ok-ng)) > 1 {
		mid := (ok + ng) / 2
		if (*srs)[mid].EpochSecond == since {
			ok = mid
			break
		}
		if (*srs)[mid].EpochSecond > since {
			ng = mid
		} else {
			ok = mid
		}
	}

	ok += 1
	sum := 0
	for _, sr := range (*srs)[ok:] {
		if since <= sr.EpochSecond && sr.EpochSecond <= until {
			sum++
		}
	}

	return sum, nil
}
