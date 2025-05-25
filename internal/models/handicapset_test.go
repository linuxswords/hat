package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHandicapEntryByBowClass(t *testing.T) {
	handicapSet := HandicapSet{
		HandicapEntries: []HandicapEntry{
			{BowClassID: 1, Value: 10.0},
			{BowClassID: 2, Value: 20.0},
			{BowClassID: 3, Value: 30.0},
		},
	}

	tests := []struct {
		bowClassID uint
		expected   *HandicapEntry
	}{
		{bowClassID: 1, expected: &HandicapEntry{BowClassID: 1, Value: 10.0}},
		{bowClassID: 2, expected: &HandicapEntry{BowClassID: 2, Value: 20.0}},
		{bowClassID: 4, expected: nil},
	}

	for _, test := range tests {
		result := handicapSet.GetHandicapEntryByBowClass(test.bowClassID)
		assert.Equal(t, test.expected, result)
	}
}
