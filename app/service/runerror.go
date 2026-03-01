package service

import (
	"context"
	"errors"
	"strings"
)

// RunErrorKind classifies a query-runner error as timeout or generic validation.
type RunErrorKind int

const (
	RunErrorValidation RunErrorKind = iota
	RunErrorTimeout
)

// ClassifyRunError inspects err from queryrunner.Run and returns the kind and a user-facing message.
func ClassifyRunError(err error) (RunErrorKind, string) {
	msg := err.Error()
	if errors.Is(err, context.DeadlineExceeded) ||
		strings.Contains(msg, "query execution timeout") ||
		strings.Contains(msg, "context deadline exceeded") ||
		strings.Contains(msg, "deadline exceeded") {
		return RunErrorTimeout, "Query timed out. Try a simpler query or reduce the amount of data."
	}
	return RunErrorValidation, msg
}
