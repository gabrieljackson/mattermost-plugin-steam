package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (p *Plugin) makeSteamAPICall(userKey, endpoint string) ([]byte, error) {
	userInfo, err := p.getSteamUserInfoByKey(userKey)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.steampowered.com/%s/?key=%s&steamid=%s&include_appinfo=true&format=json", endpoint, userInfo.APIToken, userInfo.SteamID)

	return steamAPICall(url)
}

func (p *Plugin) makeSteamAPICallSteamIDs(userKey, endpoint string) ([]byte, error) {
	userInfo, err := p.getSteamUserInfoByKey(userKey)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.steampowered.com/%s/?key=%s&steamids=%s&format=json", endpoint, userInfo.APIToken, userInfo.SteamID)

	return steamAPICall(url)
}

func steamAPICall(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
