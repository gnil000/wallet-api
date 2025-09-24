package utils

import "context"

type key = str
type str struct {
	val string
}

func ContextApp(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(key{"app"}).(str)
	if ok {
		return val.val, ok
	}
	return "", ok
}
