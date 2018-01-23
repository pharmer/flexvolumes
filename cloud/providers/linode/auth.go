package linode

import (
	"fmt"
	"os"
	"strings"

	. "github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/util"
	"github.com/taoh/linodego"
)

const (
	tokenEnv = "LINODE_API_KEY"
	apiToken = "token"
)

type TokenSource struct {
	ApiToken string `json:"token"`
}

func getCredential() (*TokenSource, error) {
	if t, err := util.ReadSecretKeyFromFile(SecretDefaultLocation, apiToken); err == nil {
		return &TokenSource{
			ApiToken: t,
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
			ApiToken: strings.TrimSpace(t),
		}, nil
	}

	cred, err := util.ReadCredentialFromFile(CredentialDefaultLocation, &TokenSource{})
	if err != nil {
		return nil, err
	}
	tokenSource := cred.(*TokenSource)
	if tokenSource.ApiToken != "" {
		return tokenSource, nil
	}

	return nil, fmt.Errorf("no credential provided for packet")
}

func (t *TokenSource) getClient() *linodego.Client {
	return linodego.NewClient(t.ApiToken, nil)
}
