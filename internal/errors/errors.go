// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package errors

import (
	stderrors "errors"
	"fmt"
)

const (
	ExitSuccess      = 0
	ExitInternal     = 1
	ExitInvalidInput = 2
	ExitUpstream     = 3
)

type Kind string

const (
	Internal       Kind = "internal"
	InvalidInput   Kind = "invalid_input"
	Upstream       Kind = "upstream"
	RateLimited    Kind = "rate_limited"
	NetworkTimeout Kind = "network_timeout"
)

type Error struct {
	Kind Kind
	Msg  string
	Err  error
}

func New(kind Kind, msg string) *Error {
	return &Error{Kind: kind, Msg: msg}
}

func Wrap(kind Kind, msg string, err error) *Error {
	return &Error{Kind: kind, Msg: msg, Err: err}
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Kind, e.Msg)
	}
	return fmt.Sprintf("%s: %s: %v", e.Kind, e.Msg, e.Err)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func KindOf(err error) Kind {
	var appErr *Error
	if stderrors.As(err, &appErr) {
		return appErr.Kind
	}
	return Internal
}

func ExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}
	switch KindOf(err) {
	case InvalidInput:
		return ExitInvalidInput
	case Upstream, RateLimited, NetworkTimeout:
		return ExitUpstream
	default:
		return ExitInternal
	}
}
