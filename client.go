package trelgo

import (
	"fmt"
	"io/ioutil"
	"net/http"

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
	// TODO check to see if expiration is valid
	c.AdditionalAuthorizationUrlParams["expiration"] = expiration
	// TODO check to see if scope is valid
	c.AdditionalAuthorizationUrlParams["scope"] = scope

	return c, nil
}
