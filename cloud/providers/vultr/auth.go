package vultr

import (
	"fmt"
	"os"
	"strings"

	vultr "github.com/JamesClonk/vultr/lib"
	. "github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/util"
)

const (
	tokenEnv = "VULTR_TOKEN"
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

	return nil, fmt.Errorf("no credential provided for vultr")
}

func (t *TokenSource) getClient() *vultr.Client {
	return vultr.NewClient(t.AccessToken, nil)
}
