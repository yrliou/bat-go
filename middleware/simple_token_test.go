package middleware

import (
	"strings"
	"testing"
)

func Test_IsSimpleTokenValid(t *testing.T) {
	tokenList := strings.Split("", ",")
	if IsSimpleTokenValid(tokenList, "") {
		t.Error("Expected empty token to always be invalid")
	}

	if IsSimpleTokenValid([]string{}, "") {
		t.Error("Expected empty token list to always to be invalid")
	}

	tokenList = strings.Split("FOO", ",")
	if !IsSimpleTokenValid(tokenList, "FOO") {
		t.Error("Expected single token to be valid")
	}
	if IsSimpleTokenValid(tokenList, "BAR") {
		t.Error("Expected wrong token to be invalid")
	}

	tokenList = strings.Split("FOO ", ",")
	if IsSimpleTokenValid(tokenList, "FOO") {
		t.Error("Expected single token to be invalid if list is malformed")
	}

	tokenList = strings.Split("FOO,BAR", ",")
	if !IsSimpleTokenValid(tokenList, "FOO") {
		t.Error("Expected multiple tokens to be valid")
	}
	if !IsSimpleTokenValid(tokenList, "BAR") {
		t.Error("Expected multiple tokens to be valid")
	}
	if IsSimpleTokenValid(tokenList, "XXX") {
		t.Error("Expected wrong tokens to be invalid")
	}
}
