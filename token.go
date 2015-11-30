package trelgo

import (
	"fmt"
	"net/http"

	"github.com/mrjones/oauth"
)

var (
	requestTokenCache *oauth.RequestToken
	accessTokenCache  *oauth.AccessToken
	tokenUrl          string
)

func HttpGetToken(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("http://%s%s", r.Host, tokenUrl)
	token, requestUrl, err := consumerCache.GetRequestTokenAndUrl(url)
	if err != nil {
		fmt.Fprintf(w, "HttpGetToken problem: %v", err)
	} else {
		requestTokenCache = token
		http.Redirect(w, r, requestUrl, http.StatusTemporaryRedirect)
	}
}

func HttpSaveToken(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	verificationCode := values.Get("oauth_verifier")
	token := values.Get("oauth_token")
	if requestTokenCache.Token != token {
		fmt.Fprintf(w, "HttpSaveToken request token did not match")
		return
	}

	accessToken, err := consumerCache.AuthorizeToken(requestTokenCache, verificationCode)
	if err != nil {
		fmt.Fprintf(w, "HttpSaveToken: %v", err)
	} else {
		accessTokenCache = accessToken
		fmt.Fprintf(w, "Trello Token Saved")
	}
}

func GetTokenInstructions(consumer *oauth.Consumer) (string, error) {
	if consumer == nil || consumerCache == nil {
		return "", fmt.Errorf("GetToken, consumer and consumerCache do not exist.")
	} else if consumer == nil {
		consumer = consumerCache
	}
	token, url, err := consumerCache.GetRequestTokenAndUrl("")
	if err != nil {
		return "", err
	}
	requestTokenCache = token
	instructions := "Go to: " + url + `
Grant access and get back a verification code.
Enter the verification code here:`
	return instructions, nil
}

func ConsoleScanForVerificationCode() (verificationCode string, err error) {
	_, err = fmt.Scanln(&verificationCode)
	return
}

func SaveToken(verificationCode string) (*oauth.AccessToken, error) {
	if consumerCache == nil {
		return nil, fmt.Errorf("SaveToken, no consumerCache found")
	}
	if requestTokenCache == nil {
		return nil, fmt.Errorf("SaveToken, no requestTokenCache found")
	}
	accessToken, err := consumerCache.AuthorizeToken(requestTokenCache, verificationCode)
	if err != nil {
		return nil, err
	}
	accessTokenCache = accessToken
	return accessToken, nil
}
