package trelgo

import (
	"flag"
	"os"
)

type Config struct {
	ConsumerKey       string
	ConsumerSecret    string
	RequestTokenUrl   string
	AuthorizeTokenUrl string
	AccessTokenUrl    string
}

func getConfig() Config {
	var (
		flagKey *string = flag.String(
			"consumerkey",
			"",
			"Consumer Key from Trello. See: https://trello.com/1/appKey/generate",
		)
		flagSecret *string = flag.String(
			"consumersecret",
			"",
			"Consumer Secret from Trello. See: https://trello.com/1/appKey/generate",
		)
		envPrefix *string = flag.String(
			"envprefix",
			"TRELLO",
			"Environment prefix for KEY and SECRET if the key and secret are not set",
		)
		requestTokenUrl   string = "https://trello.com/1/OAuthGetRequestToken"
		authorizeTokenUrl string = "https://trello.com/1/OAuthAuthorizeToken"
		accessTokenUrl    string = "https://trello.com/1/OAuthGetAccessToken"
		consumerKey       string
		consumerSecret    string
	)

	flag.Parse()

	if *flagKey == "" {
		consumerKey = os.Getenv(*envPrefix + "_KEY")
	}
	if *flagSecret == "" {
		consumerSecret = os.Getenv(*envPrefix + "_SECRET")
	}

	return Config{
		ConsumerKey:       consumerKey,
		ConsumerSecret:    consumerSecret,
		RequestTokenUrl:   requestTokenUrl,
		AuthorizeTokenUrl: authorizeTokenUrl,
		AccessTokenUrl:    accessTokenUrl,
	}
}
