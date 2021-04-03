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
	since, until := SinceUntilDate()

	userACcount, err := generateRanking(since.UnixNano(), until.UnixNano())
	if err != nil {
		log.Println(err)
	}

	sinceDay := since.Format("01/02(Mon)")
	untilDay := until.Format("01/02(Mon)")

	msg := fmt.Sprintf("%s ~ %s にかけての AC 数ランキングを発表するっぴ！\n\n", sinceDay, untilDay)
	for rank, user := range *userACcount {
		msg += fmt.Sprintf("\t%d位: <@%s> AC count: %d(%d)\n", rank+1, user.member.SlackID, user.ACcount, user.UniqueACcount)
		msg += fmt.Sprintf("\t%s\n", func() string {
			if user.ACcount < 10 {
				return "は？精進しろ"
			} else if 10 <= user.ACcount && user.ACcount < 15 {
				return "その調子っぴ！"
			} else {
				return "静大の誇りっぴ〜！"
			}
		}())
		msg += "\n"
	}

	if _, _, err := api.PostMessage(channelID, slack.MsgOptionText(msg, false)); err != nil {
		return err
	}

	return nil
}

func generateRanking(since, until int64) (*[]UserACcount, error) {
	list := []UserACcount{}

	for _, user := range *members {
		log.Println("Getting " + user.AtCoderID + "'s record...")
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
		if since <= sr.EpochSecond && sr.EpochSecond <= until && sr.Result == "AC" {
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

// 日曜日と土曜日の unixtime を返す
func SinceUntilDate() (time.Time, time.Time) {
	now := time.Now()

	sundayYear := now.Year()
	sundayMonth := now.Month()
	sunday := now.Day() - int(now.Weekday())
	if sunday < 1 {
		sundayMonth = calcMonth(now.Month(), -1)
		if now.Month() == time.January {
			sundayYear--
		}
		sunday += getDaysOfMonth(sundayYear, sundayMonth)
	}

	saturdayYear := now.Year()
	saturdayMonth := now.Month()
	saturday := now.Day() + 6 - int(now.Weekday())
	if getDaysOfMonth(now.Year(), now.Month()) < saturday {
		saturdayMonth = calcMonth(now.Month(), 1)
		if now.Month() == time.December {
			saturdayYear++
		}
		saturday -= getDaysOfMonth(saturdayYear, saturdayMonth)
	}

	sundayTime := time.Date(sundayYear, sundayMonth, sunday, 0, 0, 0, 0, time.Local)
	saturdayTime := time.Date(saturdayYear, saturdayMonth, saturday, 0, 0, 0, 0, time.Local)

	return sundayTime, saturdayTime
}

func getDaysOfMonth(year int, month time.Month) int {
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return int(lastOfMonth.Sub(firstOfMonth).Hours() / 24.0)
}

func calcMonth(currentMonth time.Month, x int) time.Month {
	if int(currentMonth)+x < 1 {
		return time.Month(int(currentMonth) - x%12)
	}

	if 12 < int(currentMonth)+x {
		return time.Month(int(currentMonth) - x%12)
	}

	return time.Month(int(currentMonth) + x)
}
