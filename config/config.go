package config

import (
	"encoding/json"
	"os"
)

type (
	Config struct {
		Hs512SecretKey  string `json:"hs512SecretKey"`
		Rs256PrivateKey string `json:"rs256PrivateKey"`
		Rs256PublicKey  string `json:"rs256PublicKey"`
		Gmail           Gmail  `json:"gmail"`
	}

	Gmail struct {
		Host          string `json:"host"`
		SenderName    string `json:"senderName"`
		SenderEmail   string `json:"senderEmail"`
		Login         string `json:"login"`
		Password      string `json:"password"`
		UrlToActivate string `json:"urlToActivate"`
		UrlToRestore  string `json:"urlToRestore"`
	}
)

func Read(filename string) (Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer func() { _ = f.Close() }()

	var config Config
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
