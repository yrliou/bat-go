package fastlycheck

import (
	"errors"
	"net/http"

	"github.com/brave-intl/bat-go/middleware"
)

var (
	FastlyTokenListKey = "FASTLY_TOKEN_LIST"
	FastlyTokenHeader  = "fastly-token"
)

func Validate(r *http.Request) error {
	list := middleware.SplitTokenList(FastlyTokenListKey)
	header := r.Header.Get(FastlyTokenHeader)
	valid := middleware.IsSimpleTokenValid(list, header)
	if !valid {
		return errors.New("an invalid token was passed")
	}
	return nil
}
