// Этот package нужен, чтобы разбить кольцевую зависимость.
package config

import (
	"github.com/cockroachdb/pebble"
	"github.com/go-redis/redis/v8"
)

// MyConfig структурка, описывающая конфиг.
type MyConfig struct {
	Server         string `json:"server,omitempty"`
	Port           int    `json:"port,omitempty"`
	Timeout        int    `json:"timeout,omitempty"`
	Loglevel       string `json:"loglevel,omitempty"`
	Log            string `json:"log,omitempty"`
	Channel        string `json:"channel,omitempty"`
	DataDir        string `json:"datadir,omitempty"`
	Csign          string `json:"csign,omitempty"`
	ForwardsMax    int64  `json:"forwards_max,omitempty"`
	OpenWeatherMap struct {
		Enabled bool   `json:"enabled,omitempty"`
		Country bool   `json:"country,omitempty"`
		Appid   string `json:"appid,omitempty"`
	} `json:"openweathermap,omitempty"`
	Flickr struct {
		Enabled          bool   `json:"enabled,omitempty"`
		Key              string `json:"key,omitempty"`
		Secret           string `json:"secret,omitempty"`
		OAuthToken       string `json:"oauth_token,omitempty"`
		OAuthTokenSecret string `json:"oauth_token_secret,omitempty"`
	} `json:"flickr,omitempty"`
	UserAgents  []string `json:"user_agents,omitempty"`
	PcacheDB    map[string]*pebble.DB
	RedisClient *redis.Client
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
