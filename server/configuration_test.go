package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigurationIsValid(t *testing.T) {
	baseConfiguration := configuration{
		EncryptionKey:         "asdf123",
		AllowedEmailDomain:    "mattermost.com",
		SteamSummaryEnable:    false,
		SteamSummaryChannelID: "",
	}

	t.Run("valid", func(t *testing.T) {
		require.NoError(t, baseConfiguration.IsValid())
	})

	t.Run("no encryption key", func(t *testing.T) {
		config := baseConfiguration
		config.EncryptionKey = ""
		require.Error(t, config.IsValid())
	})

	t.Run("cluster alerts", func(t *testing.T) {
		config := baseConfiguration
		config.SteamSummaryEnable = true
		t.Run("no channel ID", func(t *testing.T) {
			require.Error(t, config.IsValid())
		})
		t.Run("valid", func(t *testing.T) {
			config.SteamSummaryChannelID = "channel1"
			require.NoError(t, config.IsValid())
		})
	})
}
