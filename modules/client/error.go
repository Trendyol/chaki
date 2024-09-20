package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type ErrDecoder func(context.Context, *resty.Response) error

type GenericClientError struct {
	ClientName string
	StatusCode int
	RawBody    []byte
	ParsedBody interface{}
}

func (e GenericClientError) Error() string {
	msg := fmt.Sprintf("Error on client %s (Status %d)", e.ClientName, e.StatusCode)
	if details := e.extractErrorDetails(); details != "" {
		msg += ": " + details
	}
	return msg
}

func (e GenericClientError) extractErrorDetails() string {
	var details []string

	var extract func(interface{})
	extract = func(v interface{}) {
		switch value := v.(type) {
		case string:
			details = append(details, strings.TrimSpace(value))
		case map[string]interface{}:
			for _, v := range value {
				extract(v)
			}
		case []interface{}:
			for _, v := range value {
				extract(v)
			}
		}
	}

	extract(e.ParsedBody)

	if len(details) == 0 && len(e.RawBody) > 0 {
		return strings.TrimSpace(string(e.RawBody))
	}

	return strings.Join(details, "; ")
}

func DefaultErrDecoder(name string) ErrDecoder {
	return func(_ context.Context, res *resty.Response) error {
		if res.IsSuccess() {
			return nil
		}

		apiErr := GenericClientError{
			ClientName: name,
			StatusCode: res.StatusCode(),
			RawBody:    res.Body(),
		}

		var jsonBody interface{}
		if err := json.Unmarshal(res.Body(), &jsonBody); err == nil {
			apiErr.ParsedBody = jsonBody
		} else {
			apiErr.ParsedBody = string(res.Body())
		}

		return apiErr
	}
}
