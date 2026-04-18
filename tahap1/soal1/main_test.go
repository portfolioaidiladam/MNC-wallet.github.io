package main

import (
	"strings"
	"testing"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "duplikat case-sensitive acbd di posisi 2 dan 4",
			input: "4\nabcd\nacbd\naaab\nacbd\n",
			want:  "2 4",
		},
		{
			name: "Tujuh case-insensitive di posisi 3, 5, 10",
			input: "11\n" +
				"Satu\n" +
				"Dua\n" +
				"Tujuh\n" +
				"Empat\n" +
				"TUJUH\n" +
				"Enam\n" +
				"Delapan\n" +
				"Sembilan\n" +
				"Sepuluh\n" +
				"tujuh\n" +
				"Sebelas\n",
			want: "3 5 10",
		},
		{
			name:  "semua string berbeda",
			input: "5\napel\njeruk\nmangga\npisang\nanggur\n",
			want:  "false",
		},
		{
			name:  "case-insensitive di awal, set pertama berhenti di posisi 1",
			input: "3\nHalo\nHALO\nhalo\n",
			want:  "1 2 3",
		},
		{
			name:  "panjang beda tidak match",
			input: "3\nabc\nabcd\nABC\n",
			want:  "1 3",
		},
		{
			name:  "n=1",
			input: "1\nalone\n",
			want:  "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := solve(strings.NewReader(tt.input))
			if got != tt.want {
				t.Errorf("solve() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEqualIgnoreCase(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"abc", "ABC", true},
		{"Tujuh", "TUJUH", true},
		{"tujuh", "TuJuH", true},
		{"abc", "abcd", false},
		{"abc", "abd", false},
		{"", "", true},
		{"a1b", "A1B", true},
		{"a1b", "A2B", false},
	}
	for _, tt := range tests {
		got := equalIgnoreCase(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("equalIgnoreCase(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
