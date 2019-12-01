package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
)

func (p *Plugin) runSettingsCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	if len(args) == 0 {
		return nil, true, errors.New("must provide a setting")
	}
	if len(args) == 1 {
		return nil, true, errors.New("must provide setting value")
	}
	setting := args[0]
	value := args[1]

	switch setting {
	case "show-profile":
		var shown bool

		switch value {
		case "true":
			shown = true
		case "false":
			shown = false
		default:
			return nil, true, fmt.Errorf("%s is not a valid 'show-profile' setting, must be 'true' or 'false'", value)
		}

		userInfo, err := p.getSteamUserInfoByID(extra.UserId)
		if err != nil {
			return nil, true, err
		}
		if userInfo.Settings.ShowProfile == shown {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Setting %s is already %s", setting, value)), false, nil
		}
		userInfo.Settings.ShowProfile = shown
		err = p.storeSteamUser(userInfo)
		if err != nil {
			return nil, true, err
		}

	default:
		return nil, true, fmt.Errorf("%s is not a valid setting, must be 'show-profile'", value)
	}

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Setting %s updated to %s", setting, value)), false, nil
}
