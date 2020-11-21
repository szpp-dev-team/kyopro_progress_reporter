package atcoder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SubmissionResult struct {
	ID            uint32  `json:"id"`
	EpochSecond   int64   `json:"epoch_second"`
	ProblemID     string  `json:"problem_id"`
	ContentID     string  `json:"contest_id"`
	UserID        string  `json:"user_id"`
	Language      string  `json:"language"`
	Point         float64 `json:"point"`
	Length        int     `json:"length"`
	Result        string  `json:"result"`
	ExecutionTIme int     `json:"execution_time"`
}

var SubmissionResultURL = "https://kenkoooo.com/atcoder/atcoder-api/results?user=%s"

func GetSubmissionResult(userid string) (*[]SubmissionResult, error) {
	resp, err := http.Get(fmt.Sprintf(SubmissionResultURL, userid))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	submissionResult := []SubmissionResult{}
	if err := json.Unmarshal(b, &submissionResult); err != nil {
		return nil, err
	}

	return &submissionResult, nil
}

var SubmissionResultAtTimeURL = "https://kenkoooo.com/atcoder/atcoder-api/v3/from/%d"

func GetSubmissionResultAtTime(epochSecond int64) (*[]SubmissionResult, error) {
	resp, err := http.Get(fmt.Sprintf(SubmissionResultURL, epochSecond))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	buf := []SubmissionResult{}
	if err := json.Unmarshal(b, &buf); err != nil {
		return nil, err
	}

	return &buf, nil
}

var ProfileURL = "https://atcoder.jp/users/%s"

func UserExists(userid string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf(ProfileURL, userid))
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}
