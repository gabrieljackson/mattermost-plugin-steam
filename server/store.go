package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"

	"github.com/pkg/errors"
)

const (
	// StoreSteamRetries is the number of retries to use when storing steam data
	// that fails on a race.
	StoreSteamRetries = 3

	// SteamUserKey is the store suffix for a Steam profile.
	SteamUserKey = "_steam_user"

	steamAPIGetOwnedGames       = "IPlayerService/GetOwnedGames/v0001"
	steamAPIRecentlyPlayedGames = "IPlayerService/GetRecentlyPlayedGames/v0001"
	steamAPIGetSchemaForGame    = "IPlayerService/GetSchemaForGame/v0001"
)

// SteamUserInfo is the Steam profile information stored in the database.
type SteamUserInfo struct {
	MattermostUserID string        `json:"mattermost_user_id"`
	SteamID          string        `json:"steam_id"`
	APIToken         string        `json:"api_token"`
	Settings         *UserSettings `json:"user_settings"`
}

// UserSettings are user-specific settings that they can control.
type UserSettings struct {
	ShowProfile bool `json:"show_profile"`
}

func (p *Plugin) storeSteamUser(info *SteamUserInfo) error {
	config := p.getConfiguration()

	encryptedToken, err := encrypt([]byte(config.EncryptionKey), info.APIToken)
	if err != nil {
		return err
	}

	info.APIToken = encryptedToken

	jsonInfo, err := json.Marshal(info)
	if err != nil {
		return errors.Wrap(err, "unable to marshal user info")
	}

	appErr := p.API.KVSet(info.MattermostUserID+SteamUserKey, jsonInfo)
	if appErr != nil {
		return errors.Wrap(appErr, "unable to store user info in database")
	}

	return nil
}

func (p *Plugin) getSteamUserInfoByKey(key string) (*SteamUserInfo, error) {
	config := p.getConfiguration()

	var userInfo SteamUserInfo

	infoBytes, appErr := p.API.KVGet(key)
	if appErr != nil || infoBytes == nil {
		return nil, errors.New("unable to find steam user")
	}
	err := json.Unmarshal(infoBytes, &userInfo)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse user info")
	}

	unencryptedToken, err := decrypt([]byte(config.EncryptionKey), userInfo.APIToken)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decrypt steam api token")
	}

	userInfo.APIToken = unencryptedToken

	return &userInfo, nil
}

func (p *Plugin) getSteamUserInfoByID(userID string) (*SteamUserInfo, error) {
	return p.getSteamUserInfoByKey(userID + SteamUserKey)
}

func (p *Plugin) deleteSteamUser(userID string) error {
	_, err := p.getSteamUserInfoByID(userID)
	if err != nil {
		return err
	}

	appErr := p.API.KVDelete(userID + SteamUserKey)
	if appErr != nil {
		return errors.Wrap(appErr, "unable to delete user info in database")
	}

	return nil
}

func removeNonPlayerKVKeys(keys []string) []string {
	var cleanedKeys []string

	for _, key := range keys {
		if !strings.Contains(key, SteamUserKey) {
			continue
		}

		cleanedKeys = append(cleanedKeys, key)
	}

	return cleanedKeys
}

func (p *Plugin) getSteamInfoForUser(userID string) (*Player, error) {
	result, err := p.makeSteamAPICallSteamIDs(userID+SteamUserKey, "ISteamUser/GetPlayerSummaries/v0002")
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	var playerListResponse *PlayersListResponse
	err = json.Unmarshal(result, &playerListResponse)
	if err != nil {
		return nil, err
	}

	return &playerListResponse.Response.Players[0], nil
}

func encrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	msg := pad([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	finalMsg := base64.URLEncoding.EncodeToString(ciphertext)
	return finalMsg, nil
}

func decrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multiple of decoded message length")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := unpad(msg)
	if err != nil {
		return "", err
	}

	return string(unpadMsg), nil
}

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}
