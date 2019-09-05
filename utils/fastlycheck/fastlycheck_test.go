package fastlycheck

import (
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/satori/go.uuid"
)

func Test_Validate(t *testing.T) {
	alreadySetEnv := os.Getenv(FastlyTokenListKey)
	token1 := uuid.NewV4().String()
	token2 := uuid.NewV4().String()
	token3 := uuid.NewV4().String()
	originTokenList := strings.Join([]string{token1, token2, token3}, ",")
	os.Setenv(FastlyTokenListKey, originTokenList)

	finish := func(err bool, message string) {
		os.Setenv(FastlyTokenListKey, alreadySetEnv)
		if err {
			t.Error(message)
		}
	}

	req := httptest.NewRequest("GET", "/anything", nil)

	setHeader := func(value string) {
		req.Header.Set(FastlyTokenHeader, value)
	}

	if Validate(req) == nil {
		finish(true, "a token should not be set")
	}
	setHeader("")
	if Validate(req) == nil {
		finish(true, "an empty token is not allowed")
	}
	setHeader(uuid.NewV4().String())
	if Validate(req) == nil {
		finish(true, "incorrect tokens will not be allowed")
	}
	setHeader(fmt.Sprintf("Bearer %s", token1))
	if Validate(req) == nil {
		finish(true, "bearer tokens are not used")
	}
	setHeader(token1)
	if Validate(req) != nil {
		finish(true, "matching tokens should pass")
	}
	finish(false, "")
}
