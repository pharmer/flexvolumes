package util

import (
	"io/ioutil"
	"encoding/json"
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
