package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
)

const connectMessage = `Usage: |/steam connect [steam_ID] [steam_api_key]|

 - Obtain your Steam ID by viewing your Steam profile. The ID will be shown in your profile URL.
 - Obtain your API key by using the following link: https://steamcommunity.com/dev/apikey
`

func getConnectMessage() string {
	return strings.Replace(connectMessage, "|", "`", -1)
}

func (p *Plugin) runConnectCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	if len(args) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getConnectMessage()), false, nil
	}
	steamID := args[0]
	apiKey := args[1]

	// Check that the values are valid.
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s&format=json", apiKey, steamID)
	result, err := steamAPICall(url)
	if err != nil {
		return nil, true, errors.Wrap(err, "Invalid Steam credentials")
	}

	var playerListResponse *PlayersListResponse
	err = json.Unmarshal(result, &playerListResponse)
	if err != nil {
		return nil, true, errors.Wrap(err, "Invalid Steam credentials")
	}

	steamUser := &SteamUserInfo{
		MattermostUserID: extra.UserId,
		SteamID:          steamID,
		APIToken:         apiKey,
		Settings: &UserSettings{
			ShowProfile: false,
		},
	}

	err = p.storeSteamUser(steamUser)
	if err != nil {
		return nil, false, err
	}

	msg := "Steam account successfully connected!\n\n" +
		"Your profile is hidden by default. " +
		"Run `/steam settings show-profile true` to display your Steam " +
		"profile in your Mattermost Profile"

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, msg), false, nil
}

func (p *Plugin) runDisconnectCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	err := p.deleteSteamUser(extra.UserId)
	if err != nil {
		return nil, false, err
	}

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Steam account successfully disconnected."), false, nil
}
