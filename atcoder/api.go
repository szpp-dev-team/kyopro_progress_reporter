package atcoder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
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

func GetUniqueAC(userid string) (*[]SubmissionResult, error) {
	submissionResult, err := GetSubmissionResult(userid)
	if err != nil {
		return nil, err
	}

	sort.Slice(*submissionResult, func(i, j int) bool {
		return (*submissionResult)[i].ID < (*submissionResult)[j].ID
	})

	ACMap := map[string]bool{}
	UniqueACs := []SubmissionResult{}
	for _, elm := range *submissionResult {
		if elm.Result != "AC" || ACMap[elm.ProblemID] {
			continue
		}

		UniqueACs = append(UniqueACs, elm)
		ACMap[elm.ProblemID] = true
	}

	return &UniqueACs, nil
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
