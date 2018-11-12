package common

import (
	"fmt"
	"testing"
)

func TestGenUUID(t *testing.T) {
	string, _ := GenUUID()
	s, _ := GenUUID32()
	fmt.Println(string, s)
	fmt.Println(len(string), len(s))

}
