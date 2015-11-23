package trelgo

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mrjones/oauth"
)

var (
	tokens = make(map[string]*oauth.RequestToken)
	c      *oauth.Consumer
)

func TokenSetup(w http.ResponseWriter, r *http.Request) {
	tokenUrl := fmt.Sprintf("http://%s/maketoken", r.Host)
	token, requestUrl, err := c.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
		log.Fatal(err)
	}
	tokens[token.Token] = token
	http.Redirect(w, r, requestUrl, http.StatusTemporaryRedirect)
}

func TokenSave(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	verificationCode := values.Get("oauth_verifier")
	tokenKey := values.Get("oauth_token")

	accessToken, err := c.AuthorizeToken(tokens[tokenKey], verificationCode)
	if err != nil {
		log.Fatal(err)
	}
	err = overwriteFile(tokenLocation, []byte(accessTokenToString(accessToken)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Trello Token Saved")
}
