package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHandycapEntryByBowClass(t *testing.T) {
	handycapSet := HandycapSet{
		HandycapEntries: []HandycapEntry{
			{BowClassID: 1, Value: 10.0},
			{BowClassID: 2, Value: 20.0},
			{BowClassID: 3, Value: 30.0},
		},
	}

	tests := []struct {
		bowClassID uint
		expected   *HandycapEntry
	}{
		{bowClassID: 1, expected: &HandycapEntry{BowClassID: 1, Value: 10.0}},
		{bowClassID: 2, expected: &HandycapEntry{BowClassID: 2, Value: 20.0}},
		{bowClassID: 4, expected: nil},
	}

	for _, test := range tests {
		result := handycapSet.GetHandycapEntryByBowClass(test.bowClassID)
		assert.Equal(t, test.expected, result)
	}
}
