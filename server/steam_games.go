package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type GamesListResponse struct {
	Response GamesList `json:"response"`
}

type GamesList struct {
	GameCount int64  `json:"game_count"`
	Games     []Game `json:"games"`
}

// Game represents a Steam game.
type Game struct {
	AppID           int64  `json:"appid"`
	Name            string `json:"name"`
	ImgLogoURL      string `json:"img_logo_url"`
	ImgIconURL      string `json:"img_icon_url"`
	Playtime        int64  `json:"playtime_forever"`
	TwoWeekPlaytime int64  `json:"playtime_2weeks"`
	WindowsPlaytime int64  `json:"playtime_windows_forever"`
	MacPlaytime     int64  `json:"playtime_mac_forever"`
	LinuxPlaytime   int64  `json:"playtime_linux_forever"`
	StoreData       GameStoreData
}

// GameStoreDataResponse is the storefront data response for a game.
type GameStoreDataResponse struct {
	Success bool          `json:"success"`
	Data    GameStoreData `json:"data"`
}

// GameStoreData is the Steam storefront data for a game.
type GameStoreData struct {
	AppType     string           `json:"type"`
	Name        string           `json:"name"`
	AgeRequired int              `json:"required_age"`
	IsFree      bool             `json:"is_free"`
	Metacritic  GameMetacritic   `json:"metacritic"`
	Categories  []GameCategories `json:"categories"`
	Genres      []GameGenres     `json:"genres"`
}

type GameMetacritic struct {
	Score int    `json:"score"`
	URL   string `json:"url"`
}

type GameCategories struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type GameGenres struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// StoreLink returns the store link for a game.
func (g *Game) StoreLink() string {
	return fmt.Sprintf("https://store.steampowered.com/app/%d", g.AppID)
}

// PopulateStoreData obtains the Steam storefront data for a game.
func (g *Game) PopulateStoreData() error {
	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%d", g.AppID)
	result, err := steamAPICall(url)
	if err != nil {
		return err
	}

	var root map[string]GameStoreDataResponse

	err = json.Unmarshal(result, &root)
	if err != nil {
		return err
	}

	g.StoreData = root[fmt.Sprintf("%d", g.AppID)].Data

	return nil
}

func (d *GameStoreData) CategoriesToString() string {
	var categories []string
	for _, category := range d.Categories {
		categories = append(categories, category.Description)
	}

	return strings.Join(categories, ", ")
}

func (d *GameStoreData) GenresToString() string {
	var genres []string
	for _, genre := range d.Genres {
		genres = append(genres, genre.Description)
	}

	return strings.Join(genres, ", ")
}

// MakeGameMapFromRawGameListResponse returns a map of games from a raw GamesListResponse.
func MakeGameMapFromRawGameListResponse(rawGameListResponse []byte) (map[int64]Game, error) {
	var gameListResponse GamesListResponse
	err := json.Unmarshal(rawGameListResponse, &gameListResponse)
	if err != nil {
		return nil, err
	}

	gameMap := make(map[int64]Game)
	for _, game := range gameListResponse.Response.Games {
		gameMap[game.AppID] = game
	}

	return gameMap, nil
}
