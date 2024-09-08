package handler

import "testing"

func TestGetCaptcha(t *testing.T) {
	err := GetCaptcha()
	if err != nil {
		panic(err)
	}
}
