package server

import "context"

type ctxKey string

const HeadersKey ctxKey = "reqHeaders"

func GetHeaders(ctx context.Context) {
}
