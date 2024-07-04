package common

import (
	"github.com/Trendyol/chaki/util/wrapper"
	"net/http"
)

type (
	RoundTripperWrapper = wrapper.Wrapper[http.RoundTripper]
)
