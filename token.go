package trelgo

import (
	"fmt"
	"net/url"

	"github.com/mrjones/oauth"
)

// TODO make this configurable
// and not just the location of saving, but the method (database hookups?)
const tokenLocation = ".secret/token.secret"

func SetupConsumer(consumerKey, consumerSecret, name, expiration, scope string) error {
	config := getConfig()
	if consumerKey != "" {
		config.ConsumerKey = consumerKey
	}
	if consumerSecret != "" {
		config.ConsumerSecret = consumerSecret
	}
	var err error
	c, err = makeConsumer(config, name, expiration, scope)
	return err
}

func accessTokenToString(t *oauth.AccessToken) string {
	return "oauth_token=" + t.Token + "&oauth_token_secret=" + t.Secret
}

func makeAccessTokenFromString(data string) (*oauth.AccessToken, error) {
	var (
		TOKEN_PARAM        = "oauth_token"
		TOKEN_SECRET_PARAM = "oauth_token_secret"
	)
	parts, err := url.ParseQuery(data)
	if err != nil {
		return nil, err
	}

	tokenParam := parts[TOKEN_PARAM]
	parts.Del(TOKEN_PARAM)
	if len(tokenParam) < 1 {
		return nil, fmt.Errorf("Missing " + TOKEN_PARAM + " in response. " +
			"Full response body: '" + data + "'")
	}
	tokenSecretParam := parts[TOKEN_SECRET_PARAM]
	parts.Del(TOKEN_SECRET_PARAM)
	if len(tokenSecretParam) < 1 {
		return nil, fmt.Errorf("Missing " + TOKEN_SECRET_PARAM + " in response." +
			"Full response body: '" + data + "'")
	}

	additionalData := parseAdditionalData(parts)

	return &oauth.AccessToken{tokenParam[0], tokenSecretParam[0], additionalData}, nil
}

func parseAdditionalData(parts url.Values) map[string]string {
	params := make(map[string]string)
	for key, value := range parts {
		if len(value) > 0 {
			params[key] = value[0]
		}
	}
	return params
}
