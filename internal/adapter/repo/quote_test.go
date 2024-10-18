package repo

import (
	"faraway/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewQuoteRepo_Ok(t *testing.T) {
	cases := []struct {
		name   string
		text   string
		quotes []domain.Quote
	}{
		{
			name:   "happy path",
			text:   `{"quotes": [{"text":"wisdom","author":"wizard"},{"text":"magic","author":"mage"}]}`,
			quotes: []domain.Quote{{Text: "wisdom", Author: "wizard"}, {Text: "magic", Author: "mage"}},
		},
		{
			name:   "single quote",
			text:   `{"quotes": [{"text":"wisdom","author":"wizard"}]}`,
			quotes: []domain.Quote{{Text: "wisdom", Author: "wizard"}},
		},
		{
			name:   "empty list",
			text:   `{"quotes": []}`,
			quotes: []domain.Quote{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			qr, err := NewQuoteRepo([]byte(c.text))

			require.NoError(t, err)
			assert.ElementsMatch(t, c.quotes, qr.quotes)
		})
	}
}
