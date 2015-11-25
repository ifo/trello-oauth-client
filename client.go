package trelgo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/mrjones/oauth"
)

func MakeClientFromFile(consumer *oauth.Consumer, path string) (*http.Client, error) {
	if consumer == nil {
		consumer = c
	}
	if path == "" {
		path = tokenLocation
	}

	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	accessToken, err := makeAccessTokenFromString(string(rawData))
	if err != nil {
		return nil, err
	}

	client, err := consumer.MakeHttpClient(accessToken)
	if err != nil {
		return nil, err
	}
	return client, nil
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
	return fmt.Errorf("validateExpiration failed, no such expiration time exists: %s", exp)
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
