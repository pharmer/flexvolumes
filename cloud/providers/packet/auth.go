package packet

import (
	"fmt"
	"os"
	"strings"

	"github.com/packethost/packngo"
	. "github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/util"
)

const (
	tokenEnv     = "PACKET_API_KEY"
	apiKey       = "apiKey"
	projectIDKey = "projectID"
)

type TokenSource struct {
	ApiKey string `json:"apiKey"`
}

func getProjectID() (string, error) {
	return util.ReadSecretKeyFromFile(SecretDefaultLocation, projectIDKey)
}

func getCredential() (*TokenSource, error) {
	if t, err := util.ReadSecretKeyFromFile(SecretDefaultLocation, apiKey); err == nil {
		return &TokenSource{
			ApiKey: t,
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
			ApiKey: strings.TrimSpace(t),
		}, nil
	}

	cred, err := util.ReadCredentialFromFile(CredentialDefaultLocation, &TokenSource{})
	if err != nil {
		return nil, err
	}
	tokenSource := cred.(*TokenSource)
	if tokenSource.ApiKey != "" {
		return tokenSource, nil
	}

	return nil, fmt.Errorf("no credential provided for packet")
}

func (t *TokenSource) getClient() *packngo.Client {
	return packngo.NewClientWithAuth("", t.ApiKey, nil)
}
