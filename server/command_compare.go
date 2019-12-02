package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
)

func (p *Plugin) runCompareGamesCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	if len(args) == 0 {
		return nil, true, errors.New("you must provide a list of usernames to compare game lists against")
	}
	if len(args) > 10 {
		return nil, true, errors.New("the compare command is currently limited to 10 users")
	}

	var userList []string
	for _, arg := range args {
		user, err := p.API.GetUserByUsername(arg)
		if err != nil {
			return nil, true, errors.Wrapf(err, "unable to get user %s", arg)
		}

		userList = append(userList, user.Id)
	}

	// Start by getting your game list.
	result, err := p.makeSteamAPICall(extra.UserId+SteamUserKey, steamAPIGetOwnedGames)
	if err != nil {
		return nil, false, err
	}

	masterList, err := MakeGameMapFromRawGameListResponse(result)
	if err != nil {
		return nil, false, err
	}

	for _, userID := range userList {
		result, err := p.makeSteamAPICall(userID+SteamUserKey, steamAPIGetOwnedGames)
		if err != nil {
			return nil, false, err
		}

		gameMap, err := MakeGameMapFromRawGameListResponse(result)
		if err != nil {
			return nil, false, err
		}

		for appID := range masterList {
			if _, ok := gameMap[appID]; !ok {
				delete(masterList, appID)
			}
		}
	}

	output := fmt.Sprintf("Games owned by you and %s\n", strings.Join(args, ", "))
	output += fmt.Sprintf("Total: %d\n", len(masterList))
	for _, game := range masterList {
		output += fmt.Sprintf(" - [%s](%s)\n", game.Name, game.StoreLink())
	}

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, output), false, nil
}
