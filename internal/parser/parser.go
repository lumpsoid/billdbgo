package parser

import (
	"billdb/internal/bill"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// Parser defines the interface for parsing URLs.
type Parser interface {
	Parse(u *url.URL) (*bill.Bill, error)
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
func GetBillParser(u *url.URL) (Parser, error) {
	if u.Hostname() == "suf.purs.gov.rs" {
		return &ParserSerbia{}, nil
	}
	return nil, NewUnimplementedError("No parser available for the given URL")
}

func DateParse(dateLayout string, dateString string) (*time.Time, error) {
	dateTime, err := time.Parse(dateLayout, dateString)
	if err != nil {
		log.Error("Error parsing date:", err)
		return nil, err
	}
	return &dateTime, nil
}
