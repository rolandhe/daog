package utils

import (
	"fmt"
	"testing"
)

func TestToHexString(t *testing.T) {

	data := []byte{0xF0,0xFE,0x8}

	s := ToHexString(data,true)

	if s != "F0FE08"{
		t.Error("error")
	}
	fmt.Println(s)
}