package trelgo

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/mrjones/oauth"
)

var (
	consumerCache *oauth.Consumer
)

func MakeClient(consumer *oauth.Consumer, accessToken *oauth.AccessToken) (*http.Client, error) {
	if accessToken == nil && accessTokenCache == nil {
		return nil, fmt.Errorf("MakeClient failed - no oauth.AccessToken")
	}
	if accessToken == nil {
		accessToken = accessTokenCache
	}
	// TODO: check to see if a consumer needs a Key and Secret to make a Client
	// if the access token already exists and is valid
	if consumer == nil && consumerCache == nil {
		return nil, fmt.Errorf("MakeClient failed - no consumer")
	}
	if consumer == nil {
		consumer = consumerCache
	}
	return consumer.MakeHttpClient(accessToken)
}

func SetupConsumer(consumerKey, consumerSecret, name, expiration, scope string) error {
	config := getConfig()
	if consumerKey != "" {
		config.ConsumerKey = consumerKey
	}
	if consumerSecret != "" {
		config.ConsumerSecret = consumerSecret
	}
	consumer, err := makeConsumer(config, name, expiration, scope)
	if err != nil {
		return err
	}
	consumerCache = consumer
	return nil
}

func makeConsumer(conf Config, name, expiration, scope string) (*oauth.Consumer, error) {
	c := oauth.NewConsumer(
		conf.ConsumerKey,
		conf.ConsumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   conf.RequestTokenUrl,
			AuthorizeTokenUrl: conf.AuthorizeTokenUrl,
			AccessTokenUrl:    conf.AccessTokenUrl,
		},
	)

	if name == "" {
		return nil, fmt.Errorf("Application name required for oauth consumer")
	}
	c.AdditionalAuthorizationUrlParams["name"] = name
	err := validateExpiration(expiration)
	if err != nil {
		return nil, err
	}
	c.AdditionalAuthorizationUrlParams["expiration"] = expiration
	err = validateScope(scope)
	if err != nil {
		return nil, err
	}
	c.AdditionalAuthorizationUrlParams["scope"] = scope

	return c, nil
}

func validateExpiration(exp string) error {
	exipirations := []string{"1hour", "1day", "30days", "never"}
	for _, e := range exipirations {
		if e == exp {
			return nil
		}
	}
	return fmt.Errorf("validateExpiration, no such expiration time: %s", exp)
}

func validateScope(scope string) error {
	scopes := strings.Split(scope, ",")
	sort.Strings(scopes)
	validScopeLists := [][]string{
		[]string{"read"},
		[]string{"read", "write"},
		[]string{"account", "read"},
		[]string{"account", "read", "write"}}

	for _, s := range validScopeLists {
		if compareStringSlices(s, scopes) {
			return nil
		}
	}
	return fmt.Errorf("validateScope failed, no such scope exists: %s", scope)
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
