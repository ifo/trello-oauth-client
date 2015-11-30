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
