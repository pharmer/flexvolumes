package util

import (
	"io/ioutil"
	"encoding/json"
	"strings"
)

func ReadCredentialFromFile(file string, cred interface{}) (interface{}, error) {
	c, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(c, cred)
	if err != nil {
		return nil, err
	}
	return cred, nil

}

func ReadSecretKeyFromFile(location, key string) (string, error)  {
	file := strings.Join([]string{location, key}, "/")

	k, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(k)), nil
}