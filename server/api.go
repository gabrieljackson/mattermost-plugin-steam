package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mattermost/mattermost-server/plugin"
	"github.com/pkg/errors"
)

//PlayersListResponse is an API response for a list of Steam players.
type PlayersListResponse struct {
	Response PlayersList `json:"response"`
}

// PlayersList is a list of Steam players.
type PlayersList struct {
	Players []Player `json:"players"`
}

// Player is Steam player information.
type Player struct {
	SteamID     string `json:"steamid"`
	PersonaName string `json:"personaname"`
	ProfileURL  string `json:"profileurl"`
	Avatar      string `json:"avatar"`
}

// SteamUserInfoRequest is the request type to obtain steam info for a given user.
type SteamUserInfoRequest struct {
	UserID string `json:"user_id"`
}

// ServeHTTP handles HTTP requests to the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()

	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch path := r.URL.Path; path {
	case "/profile.png":
		p.handleProfileImage(w, r)
	case "/api/v1/userinfo":
		p.handleUserInfo(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleUserInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var userInfoRequest SteamUserInfoRequest
	err := json.NewDecoder(r.Body).Decode(&userInfoRequest)
	if err != nil || userInfoRequest.UserID == "" {
		if err != nil {
			p.API.LogError(errors.Wrap(err, "Unable to decode steam user request").Error())
		}

		http.Error(w, "Please provide a JSON object with a non-blank user_id field", http.StatusBadRequest)
		return
	}

	userInfo, err := p.getSteamUserInfoByID(userInfoRequest.UserID)
	if err != nil {
		p.API.LogError(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !userInfo.Settings.ShowProfile {
		w.WriteHeader(http.StatusOK)
		return
	}

	steamUserInfo, err := p.getSteamInfoForUser(userInfoRequest.UserID)
	if err != nil {
		p.API.LogError(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(steamUserInfo)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (p *Plugin) handleProfileImage(w http.ResponseWriter, r *http.Request) {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.API.LogError("Unable to get bundle path, err=" + err.Error())
		return
	}

	img, err := os.Open(filepath.Join(bundlePath, "assets", "profile.png"))
	if err != nil {
		http.NotFound(w, r)
		p.API.LogError("Unable to read profile image, err=" + err.Error())
		return
	}
	defer img.Close()

	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, img)
}
