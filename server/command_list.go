package main

import (
	"encoding/json"
	"fmt"

	"github.com/mattermost/mattermost-server/model"
)

func (p *Plugin) runListGamesCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	result, err := p.makeSteamAPICall(extra.UserId+SteamUserKey, STEAM_API_GETOWNEDGAMES)
	if err != nil {
		return nil, false, err
	}

	var gameListResponse GamesListResponse
	err = json.Unmarshal(result, &gameListResponse)
	if err != nil {
		return nil, false, err
	}

	var output string
	for _, game := range gameListResponse.Response.Games {
		output += fmt.Sprintf("- [%s](%s)\n", game.Name, game.StoreLink())
	}

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, output), false, nil
}
