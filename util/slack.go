package util

import "github.com/slack-go/slack"

// slack id -> slack name
var SlackIdNameMap map[string]string = map[string]string{}

func ConvertIdToName(api *slack.Client) error {
	list, err := api.GetUsers()
	if err != nil {
		return err
	}

	for _, elem := range list {
		SlackIdNameMap[elem.ID] = elem.Name
	}

	return nil
}
