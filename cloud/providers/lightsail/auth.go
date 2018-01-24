package lightsail

import (
	"fmt"
	"os"

	_aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
	. "github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/util"
)

const (
	accessKeyIDEnv     = "AWS_ACCESS_KEY_ID"
	secretAccessKeyEnv = "AWS_SECRET_ACCESS_KEY"
	accessKeyID        = "accessKeyID"
	secretAccessKey    = "secretAccessKey"
)

type TokenSource struct {
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
}

func getCredential() (*TokenSource, error) {
	tkn := &TokenSource{}
	if k, err := util.ReadSecretKeyFromFile(SecretDefaultLocation, accessKeyID); err == nil {
		tkn.AccessKeyID = k
	}

	if s, err := util.ReadSecretKeyFromFile(SecretDefaultLocation, secretAccessKey); err == nil {
		tkn.SecretAccessKey = s
		return tkn, nil
	}

	if f, ok := os.LookupEnv(CredentialFileEnv); ok && f != "" {
		cred, err := util.ReadCredentialFromFile(f, &TokenSource{})
		if err != nil {
			return nil, err
		}
		return cred.(*TokenSource), nil
	}

	if t, ok := os.LookupEnv(accessKeyIDEnv); ok && t != "" {
		tkn.AccessKeyID = t
	}
	if t, ok := os.LookupEnv(secretAccessKeyEnv); ok && t != "" {
		tkn.SecretAccessKey = t
		return tkn, nil
	}

	cred, err := util.ReadCredentialFromFile(CredentialDefaultLocation, &TokenSource{})
	if err != nil {
		return nil, err
	}
	tokenSource := cred.(*TokenSource)
	if tokenSource.AccessKeyID != "" && tokenSource.SecretAccessKey != "" {
		return tokenSource, nil
	}

	return nil, fmt.Errorf("no credential provided for digitalocean")
}

func (t *TokenSource) getClient() (*lightsail.Lightsail, error) {
	region, err := getRegion()
	if err != nil {
		return nil, err
	}

	config := &_aws.Config{
		Region:      _aws.String(region),
		Credentials: credentials.NewStaticCredentials(t.AccessKeyID, t.SecretAccessKey, ""),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	return lightsail.New(sess), nil
}
