package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/earlgray283/kyopro_progress_reporter/util"
)

type Member struct {
	SlackID     string `json:"slack_id"`
	AtCoderID   string `json:"atcoder_id"`
	CafeCoderID string `json:"cafecoder_id"`
	TwitterID   string `json:"twitter_id"`
}

func NewMemberFromJSON() (*[]Member, error) {
	b, err := ioutil.ReadFile("members.json")
	if err != nil {
		return nil, err
	}

	member := []Member{}
	if err := json.Unmarshal(b, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

func AddMemberToJSON() error {
	f, err := ioutil.ReadFile("member.json")
	if err != nil {
		return err
	}

	_ = os.Remove("member.json.bk")
	if err := ioutil.WriteFile("member.json.bk", f, 0666); err != nil {
		return err
	}

	b, err := json.Marshal(*members)
	if err != nil {
		return err
	}

	_ = os.Remove("member.json")
	if err := ioutil.WriteFile("member.json", b, 0666); err != nil {
		return err
	}

	if err := util.DownloadFile("member.json"); err != nil {
		return err
	}
	if err := util.DownloadFile("member.json.bk"); err != nil {
		return err
	}

	return nil
}

func AddMember(jsonMsg string) error {
	member := Member{}
	if err := json.Unmarshal([]byte(jsonMsg), &member); err != nil {
		return err
	}

	*members = append(*members, member)

	if err := AddMemberToJSON(); err != nil {
		return err
	}

	return nil
}
