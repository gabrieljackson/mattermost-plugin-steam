package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
)

type recentGame struct {
	AppID    int64
	Playtime int64
}

func (p *Plugin) runListRecentGamesCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	keys, appErr := p.API.KVList(0, 1000)
	if appErr != nil {
		return nil, false, appErr
	}
	keys = removeNonPlayerKVKeys(keys)

	var totalPlaytime int64
	gamesPlayed := make(map[int64]int64)
	gamesReference := make(map[int64]Game)

	for _, key := range keys {
		result, err := p.makeSteamAPICall(key, steamAPIRecentlyPlayedGames)
		if err != nil {
			p.API.LogError(errors.Wrapf(err, "unable to get recently-played games for %s", key).Error())
			continue
		}

		var gameListResponse GamesListResponse
		err = json.Unmarshal(result, &gameListResponse)
		if err != nil {
			return nil, false, err
		}

		for _, game := range gameListResponse.Response.Games {
			totalPlaytime += game.TwoWeekPlaytime
			gamesPlayed[game.AppID] += game.TwoWeekPlaytime
			gamesReference[game.AppID] = game
		}
	}

	var gamesPlayedSlice []recentGame
	for k, v := range gamesPlayed {
		gamesPlayedSlice = append(gamesPlayedSlice, recentGame{AppID: k, Playtime: v})
	}
	sort.Slice(gamesPlayedSlice, func(i, j int) bool { return gamesPlayedSlice[i].Playtime > gamesPlayedSlice[j].Playtime })

	output := fmt.Sprintf("Recently Played Summary for %d Players [%d minutes total]:\n\n", len(keys), totalPlaytime)
	for _, recentGame := range gamesPlayedSlice {
		game := gamesReference[recentGame.AppID]
		output += fmt.Sprintf(" - [%s](%s) [%d minutes]\n", game.Name, game.StoreLink(), recentGame.Playtime)
	}

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, output), false, nil
}
