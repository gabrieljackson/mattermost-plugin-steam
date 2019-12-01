package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const helpText = `* |/steam connect| - Connect your Mattermost account to your Steam account
* |/steam disconnect| - Disconnect your Mattermost account from your Steam account
* |/steam list| - Shows the list of games in your Steam library
* |/steam recent| - Shows recent game stats about other Steam plugin users
* |/steam compare [user1] [user2] [user3] [etc.]| - Compare owned games with one or multiple other Steam plugin users
* |/steam settings [setting] [value]| - Update your user settings
  * |setting| can be "show-profile"
  * |value| can be "true" or "false"
* |/steam info| - Shows plugin information`

func getHelp() string {
	return strings.Replace(helpText, "|", "`", -1)
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "steam",
		DisplayName:      "Steam",
		Description:      "Integration with Steam",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: connect, disconnect, recent, compare, settings, info",
		AutoCompleteHint: "[command]",
	}
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     "steam",
		IconURL:      fmt.Sprintf("/plugins/%s/profile.png", manifest.ID),
	}
}

// ExecuteCommand executes a given command and returns a command response.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	config := p.getConfiguration()

	if config.AllowedEmailDomain != "" {
		user, err := p.API.GetUser(args.UserId)
		if err != nil {
			return nil, err
		}

		if !strings.HasSuffix(user.Email, config.AllowedEmailDomain) {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Permission denied. Please talk to your system administrator to get access."), nil
		}
	}

	stringArgs := strings.Split(args.Command, " ")

	if len(stringArgs) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	}

	command := stringArgs[1]

	var handler func([]string, *model.CommandArgs) (*model.CommandResponse, bool, error)

	switch command {
	case "connect":
		handler = p.runConnectCommand
	case "disconnect":
		handler = p.runDisconnectCommand
	case "list":
		handler = p.runListGamesCommand
	case "compare":
		handler = p.runCompareGamesCommand
	case "recent":
		handler = p.runListRecentGamesCommand
	case "settings":
		handler = p.runSettingsCommand
	case "info":
		handler = p.runInfoCommand
	}

	if handler == nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	}

	resp, userError, err := handler(stringArgs[2:], args)

	if err != nil {
		p.API.LogError(err.Error())
		if userError {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("__Error: %s__\n\nRun `/steam help` for usage instructions.", err.Error())), nil
		}

		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "An unknown error occurred. Please talk to your administrator for help."), nil
	}

	return resp, nil
}

func (p *Plugin) runInfoCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	resp := fmt.Sprintf("Steam plugin version: %s, "+
		"[%s](https://github.com/gabrieljackson/mattermost-steam-plugin/commit/%s), built %s\n\n",
		manifest.Version, BuildHashShort, BuildHash, BuildDate)

	keys, appErr := p.API.KVList(0, 1000)
	if appErr != nil {
		return nil, false, appErr
	}
	keys = removeNonPlayerKVKeys(keys)

	resp += "Stats:\n"
	resp += fmt.Sprintf(" - Plugin Users: %d\n", len(keys))

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, resp), false, nil
}
