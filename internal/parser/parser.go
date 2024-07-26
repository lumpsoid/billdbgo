package parser

import (
	"billdb/internal/bill"
)

// Parser defines the interface for parsing URLs.
type Parser interface {
	Parse(u string) (*bill.Bill, error)
}

// UnimplementedError represents an unimplemented feature error.
type UnimplementedError struct {
	message string
}

// Error returns the error message.
func (e *UnimplementedError) Error() string {
	return e.message
}

// NewUnimplementedError creates a new UnimplementedError instance with the given message.
func NewUnimplementedError(message string) *UnimplementedError {
	return &UnimplementedError{message: message}
}

// GetBillParser creates a parser for a given URL.
func GetBillParser(data string) (Parser, error) {
  switch data {
  case "https://suf.purs.gov.rs":
		return &ParserSerbia{}, nil
  }
	return nil, NewUnimplementedError("No parser available for the given URL")
}
