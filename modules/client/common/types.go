package common

import (
	"net/http"

	"github.com/Trendyol/chaki/util/wrapper"
)

type (
	RoundTripperWrapper = wrapper.Wrapper[http.RoundTripper]
)
