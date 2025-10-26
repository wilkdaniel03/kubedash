package common

import (
	"errors"
	"os"
)

var serversecretkey []byte = nil

func InitSSK(secretkeyPath string) error {

	secretkeyData, err := os.ReadFile(secretkeyPath)
	if err != nil {
		return err
	}

	serversecretkey = secretkeyData

	return nil

}

func GetSSK() ([]byte, error) {

	if serversecretkey == nil {
		return nil, errors.New("could not get secret key")
	}

	return serversecretkey, nil
}
