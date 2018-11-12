package common

import "encoding/hex"

import (
	"crypto/rand"
	"strings"
)

func GenUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// TODO: verify the two lines implement RFC 4122 correctly
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return hex.EncodeToString(uuid), nil
}

func GenUUID32() (string, error) {
	uuid, err := GenUUID()
	uuid = strings.Replace(uuid, "-", "", -1)
	return uuid, err
}
