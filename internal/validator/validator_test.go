package validator

import "testing"

func TestValidateTrade(t *testing.T) {

	tests := []struct {
		name     string
		account  string
		symbol   string
		side     string
		volume   float64
		open     float64
		close    float64
		expected error
	}{
		{
			name:     "valid trade",
			account:  "acc123",
			symbol:   "EURUSD",
			side:     "buy",
			volume:   1.0,
			open:     1.1,
			close:    1.2,
			expected: nil,
		},
		{
			name:     "empty account",
			account:  "",
			symbol:   "EURUSD",
			side:     "buy",
			volume:   1.0,
			open:     1.1,
			close:    1.2,
			expected: ErrEmptyAccount,
		},
		{
			name:     "invalid symbol",
			account:  "acc123",
			symbol:   "iasdfo",
			side:     "buy",
			volume:   1.0,
			open:     1.1,
			close:    1.2,
			expected: ErrInvalidSymbol,
		},
		{
			name:     "volume must be positive",
			account:  "acc123",
			symbol:   "EURUSD",
			side:     "buy",
			volume:   0.0,
			open:     1.1,
			close:    1.2,
			expected: ErrInvalidVolume,
		},
		{
			name:     "invalid price",
			account:  "acc123",
			symbol:   "EURUSD",
			side:     "buy",
			volume:   1.0,
			open:     0.0,
			close:    1.0,
			expected: ErrInvalidPrice,
		},
		{
			name:     "invalid side",
			account:  "acc123",
			symbol:   "EURUSD",
			side:     "buyy",
			volume:   1.0,
			open:     1.1,
			close:    1.2,
			expected: ErrInvalidSide,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := ValidateTrade(tt.account, tt.symbol, tt.side, tt.volume, tt.open, tt.close)
			if err != tt.expected {
				t.Errorf("ValidateTrade() error = %v, expected %v", err, tt.expected)
			}

		})
	}

}
