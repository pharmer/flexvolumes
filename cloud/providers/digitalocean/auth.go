package digitalocean

import (
	"fmt"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	. "github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/util"
	"golang.org/x/oauth2"
)

const (
	tokenEnv = "DO_ACCESS_TOKEN"
	tokenKey = "token"
)

type TokenSource struct {
	AccessToken string `json:"token"`
}

func getCredential() (*TokenSource, error) {
	if t, err := util.ReadSecretKeyFromFile(SecretDefaultLocation, tokenKey); err == nil {
		return &TokenSource{
			AccessToken: t,
		}, nil
	}

	if f, ok := os.LookupEnv(CredentialFileEnv); ok && f != "" {
		cred, err := util.ReadCredentialFromFile(f, &TokenSource{})
		if err != nil {
			return nil, err
		}
		return cred.(*TokenSource), nil
	}

	if t, ok := os.LookupEnv(tokenEnv); ok && t != "" {
		return &TokenSource{
			AccessToken: strings.TrimSpace(t),
		}, nil
	}

	cred, err := util.ReadCredentialFromFile(CredentialDefaultLocation, &TokenSource{})
	if err != nil {
		return nil, err
	}
	tokenSource := cred.(*TokenSource)
	if tokenSource.AccessToken != "" {
		return tokenSource, nil
	}

	return nil, fmt.Errorf("no credential provided for digitalocean")
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func (t *TokenSource) getClient() *godo.Client {
	tokenSource := &TokenSource{
		AccessToken: t.AccessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	return godo.NewClient(oauthClient)
}
