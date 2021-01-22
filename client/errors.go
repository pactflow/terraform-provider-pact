package client

import (
	"errors"
	"fmt"
	"strings"
)

type APIErrorKey string
type APIErrorDescriptions []string

type APIKeyedError map[APIErrorKey]APIErrorDescriptions

// APIErrorResponse represents a body of the shape: {"errors":{"parameter": ["parameter is invaled"]}}
type APIErrorResponse struct {
	Errors       APIKeyedError   `json:"errors"`
	Reference    string          `json:"reference"`
	ErrorDetails APIErrorMessage `json:"error"`
	err          error
}

// APIErrorMessage represents are higher-level error such as for the cause of a 5xx
type APIErrorMessage struct {
	Message string `json:"message"`
}

// APIArrayErrorResponse represents a body of the shape: {"errors":["Missing required attribute 'name'"]}
type APIArrayErrorResponse struct {
	Errors       APIErrorDescriptions `json:"errors"`
	Reference    string               `json:"reference"`
	ErrorDetails APIErrorMessage      `json:"error"`
	err          error
}

type APIError interface {
	Error() string
	GetError() error
	GetReference() string
	GetMessage() string
}

// 400 bad request message
// {"errors":["Missing required attribute 'name'"]}

// 500 error message
// {"error":{"message":

func (e *APIErrorResponse) Error() string {
	errors := new(strings.Builder)
	errors.WriteString("\terror details: \n")

	if e.ErrorDetails.Message != "" {
		errors.WriteString(fmt.Sprintf("\t\tsummary: %s\n", e.ErrorDetails.Message))
	}

	if len(e.Errors) > 0 {
		for _, v := range e.Errors {
			// errors.WriteString(fmt.Sprintf("\t\t%s - %s\n", k, strings.Join(v, ",")))
			errors.WriteString(fmt.Sprintf("\t\t%s\n", strings.Join(v, "\n")))
		}
	}

	if e.Reference != "" {
		errors.WriteString(fmt.Sprintf("\t\treference: %s\n", e.Reference))
	}

	if e.err != nil {
		return fmt.Sprintf("\t%s \n\n%s", e.err, errors.String())
	}

	return errors.String()
}

func (e *APIArrayErrorResponse) Error() string {
	errors := new(strings.Builder)
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

	if e.err != nil {
		return fmt.Sprintf("\t%s\n\n%s", e.err, errors.String())
	}

	return errors.String()
}

var (
	ErrBadRequest        = errors.New("bad request")
	ErrSystemUnavailable = errors.New("system unavailable")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("access denied, check that you have access to this resource")
)
