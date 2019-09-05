package fastlycheck

import (
	"net/http"

	"github.com/brave-intl/bat-go/middleware"
)

var (
	FastlyTokenListKey = "FASTLY_TOKEN_LIST"
	FastlyTokenHeader  = "fastly-token"
)

func Validate(r *http.Request) bool {
	list := middleware.SplitTokenList(FastlyTokenListKey)
	header := r.Header.Get(FastlyTokenHeader)
	return middleware.IsSimpleTokenValid(list, header)
}
