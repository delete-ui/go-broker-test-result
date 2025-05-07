package validator

import (
	"errors"
	"regexp"
)

var (
	ErrEmptyAccount  error = errors.New("account cannot be empty")
	ErrInvalidSymbol error = errors.New("symbol must be 6 uppercase letters")
	ErrInvalidVolume error = errors.New("volume must be positive")
	ErrInvalidPrice  error = errors.New("open and close price must be positive")
	ErrInvalidSide   error = errors.New("side must be either 'buy' or 'sell'")

	symbolRegex = regexp.MustCompile("^[A-Z]{6}$")
)

func ValidateTrade(account, symbol, side string, volume, open, close float64) error {

	if account == "" {
		return ErrEmptyAccount
	}

	if !symbolRegex.MatchString(symbol) {
		return ErrInvalidSymbol
	}

	if volume <= 0 {
		return ErrInvalidVolume
	}

	if open <= 0 || close <= 0 {
		return ErrInvalidPrice
	}

	switch side {
	case "buy", "sell":
		return nil
	default:
		return ErrInvalidSide
	}

}
