package client

import (
	"errors"
	"fmt"
	"strings"
)

type apiErrorKey string
type apiErrorDescriptions []string

type APIKeyedError map[apiErrorKey]apiErrorDescriptions

// apiErrorResponse represents a body of the shape: {"errors":{"parameter": ["parameter is invaled"]}}
type apiErrorResponse struct {
	Errors       APIKeyedError   `json:"errors"`
	Reference    string          `json:"reference"`
	ErrorDetails apiErrorMessage `json:"error"`
	err          error
}

// apiErrorMessage represents are higher-level error such as for the cause of a 5xx
type apiErrorMessage struct {
	Message string `json:"message"`
}

// apiArrayErrorResponse represents a body of the shape: {"errors":["Missing required attribute 'name'"]}
type apiArrayErrorResponse struct {
	Errors       apiErrorDescriptions `json:"errors"`
	Reference    string               `json:"reference"`
	ErrorDetails apiErrorMessage      `json:"error"`
	err          error
}

type apiError interface {
	Error() string
	GetError() error
	GetReference() string
	GetMessage() string
}

func (e *apiErrorResponse) Error() string {
	errors := new(strings.Builder)
	if e.ErrorDetails.Message != "" || len(e.Errors) > 0 || e.Reference != "" {
		errors.WriteString("\terror details: \n")

		if e.ErrorDetails.Message != "" {
			errors.WriteString(fmt.Sprintf("\t\tsummary: %s\n", e.ErrorDetails.Message))
		}

		if len(e.Errors) > 0 {
			for _, v := range e.Errors {
				errors.WriteString(fmt.Sprintf("\t\t%s\n", strings.Join(v, "\n")))
			}
		}

		if e.Reference != "" {
			errors.WriteString(fmt.Sprintf("\t\treference: %s\n", e.Reference))
		}
	} else {
		errors.WriteString("\n\tplease see the log for error details\n")
	}

	if e.err != nil {
		return fmt.Sprintf("\t%s \n\n%s", e.err, errors.String())
	}

	return errors.String()
}

func (e *apiArrayErrorResponse) Error() string {
	errors := new(strings.Builder)
	if e.ErrorDetails.Message != "" || len(e.Errors) > 0 || e.Reference != "" {
		errors.WriteString("\n\terror details: \n")

		if e.ErrorDetails.Message != "" {
			errors.WriteString(fmt.Sprintf("\t\tsummary: %s\n", e.ErrorDetails.Message))
		}

		if len(e.Errors) > 0 {
			errors.WriteString(fmt.Sprintf("\t\t%s\n", strings.Join(e.Errors, "\n")))
		}

		if e.Reference != "" {
			errors.WriteString(fmt.Sprintf("\t\treference: %s\n", e.Reference))
		}
	} else {
		errors.WriteString("\n\tplease see the log for error details\n")
	}

	if e.err != nil {
		return fmt.Sprintf("\t%s\n\n%s", e.err, errors.String())
	}

	return errors.String()
}

var (
	// ErrBadRequest represents an HTTP 400 error
	ErrBadRequest = errors.New("bad request")
	// ErrSystemUnavailable represents an hTTP 5xx error
	ErrSystemUnavailable = errors.New("system unavailable")
	// ErrUnauthorized represents an HTTP 401 error
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden represents an HTTP 403 permissions issue
	ErrForbidden = errors.New("access denied, check that you have access to this resource")
)
