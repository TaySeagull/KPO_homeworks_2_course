package domain_test

import (
	"testing"

	"hw_2/internal/domain"
)

func TestNewScoreValue(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:  "valid score",
			value: 5,
		},
		{
			name:  "zero score",
			value: 0,
		},
		{
			name:    "negative score",
			value:   -1,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			score, err := domain.NewScoreValue(test.value)

			if (err != nil) != test.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if err == nil && score.Value() != test.value {
				t.Errorf("got %d, want %d", score.Value(), test.value)
			}
		})
	}
}

func TestNewMaxScore(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:  "valid maximum",
			value: 10,
		},
		{
			name:    "zero maximum",
			value:   0,
			wantErr: true,
		},
		{
			name:    "negative maximum",
			value:   -5,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			maxScore, err := domain.NewMaxScore(test.value)

			if (err != nil) != test.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if err == nil && maxScore.Value() != test.value {
				t.Errorf("got %d, want %d", maxScore.Value(), test.value)
			}
		})
	}
}

func TestNewWeight(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:  "minimum weight",
			value: 1,
		},
		{
			name:  "maximum weight",
			value: 100,
		},
		{
			name:  "valid weight",
			value: 50,
		},
		{
			name:    "zero weight",
			value:   0,
			wantErr: true,
		},
		{
			name:    "negative weight",
			value:   -10,
			wantErr: true,
		},
		{
			name:    "weight greater than 100",
			value:   101,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			weight, err := domain.NewWeight(test.value)

			if (err != nil) != test.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if err == nil && weight.Value() != test.value {
				t.Errorf("got %d, want %d", weight.Value(), test.value)
			}
		})
	}
}

func TestNewFinalScore(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:  "minimum final score",
			value: 0,
		},
		{
			name:  "maximum final score",
			value: 10,
		},
		{
			name:  "valid final score",
			value: 7,
		},
		{
			name:    "negative final score",
			value:   -1,
			wantErr: true,
		},
		{
			name:    "final score greater than 10",
			value:   11,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			score, err := domain.NewFinalScore(test.value)

			if (err != nil) != test.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if err == nil && score.Value() != test.value {
				t.Errorf("got %d, want %d", score.Value(), test.value)
			}
		})
	}
}
