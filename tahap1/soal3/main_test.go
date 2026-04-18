package main

import "testing"

// Long bracket strings dari PDF (konkat multi-line -> single string).
const (
	truePDF3 = `{{[{<[[{<{<<<[{{{[]{<{[<[[<{{[[[[[<{[{<[<<[[<<{[[{[<<<<<<<[{[{[{{<{[[<{<<<{<{[<>]}>}>>[]>}>]]}>}}]}]}]>>>>>>>]}]]}>>]]>>]>}]}>]]]]]}}>]]>]}>}}}}]>>>}>}]]>}]}}`

	truePDF4 = `[<{<{[{[{}[[<[<{{[<[<[[[<{{[<<<[[[<[<{{[<<{{<{<{<[<{[{{[{{{{[<<{{{<{[{[[[{<<<[{[<{<<<>>>}>]}]>>>}]]]}]}>}}}>>]}}}}]}}]}>]>}>}>}}>>]}}>]>]]]>>>]}}>]]]>]>]}}>]>]]]}]}>}>]`

	truePDF5 = `[[{[[<{{{{[[[<[{[[<{{{{[{[{[[[<<{<{[{<<<[[<[{[<{{[{[<[[<<[{<<[[[{<[{[[{{<<>[<<{{<<{[[[<{}{[{{{[[{{[[<[{}]>]]}}]]}}}]}>]]]}>>}}>>]>}}]]}]>}]]]>>}]>>]]>]}]}}>]}]>]]>>>}]}>}>>]]]}]}]}}}}>]]}]>]]]}}}}>]]}]]`

	falsePDF7 = `{{[{<[[{<{<<<[{{{[]{<{[<[[<{{[[[[<{[{<[<<[[<<{[[{[<<<<<<<[{[{[{{<{[[<{<<<{<{[<>]}>}>>[]>}>]]}>}}]}]}]>>>>>>]}]]}>>]]>>]>}]}>]]]]]}}>]]>]}>}}}}]>>>}>}]]>}]}}`

	falsePDF8 = `[<{<{[{[{}[[<[<{{[<[<[[[<{{[<<<[[[<[<{{[<<{{<{<{<[<{[{{[{{{{[<<{{{<{[{[[[{<<<[{[<{<<>>[]}]>>>}]]]}]}>}}}>>]}}}}]}}]}>]>}>}>}}>>]}}>]>]]]>>>]}}>]]]>]>]}}>]>]]]}]}>}>]`
)

func TestValidateBrackets_True(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"PDF TRUE 1", "{{[<>[{{}}]]}}"},
		{"PDF TRUE 2", "{<{[[{{[]<{{[{[]<>}]}}<>>}}]]}>}"},
		{"PDF TRUE 3 (long)", truePDF3},
		{"PDF TRUE 4 (long)", truePDF4},
		{"PDF TRUE 5 (long)", truePDF5},
		{"PDF TRUE 6", "[{}<>]"},
		{"rule nested penuh <[]{<>}>", "<[]{<>}>"},
		{"empty-looking simple", "[]"},
		{"semua 3 jenis", "<{[]}>"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if !ValidateBrackets(tt.in) {
				t.Errorf("ValidateBrackets(%q) = false, want true", tt.in)
			}
		})
	}
}

func TestValidateBrackets_False(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"PDF FALSE 1: ]", "]"},
		{"PDF FALSE 2: ][", "]["},
		{"PDF FALSE 3: [>]", "[>]"},
		{"PDF FALSE 4: [>", "[>"},
		{"PDF FALSE 6 (mirip TRUE 2 minus 1 kurawal)", "{<{[[{{[]<{[{[]<>}]}}<>>}}]]}>}"},
		{"PDF FALSE 7 (long)", falsePDF7},
		{"PDF FALSE 8 (long)", falsePDF8},
		{"PDF FALSE 10: [{}<[>]", "[{}<[>]"},
		{"rule silang <[>]", "<[>]"},
		{"pembuka tanpa penutup", "<{["},
		{"penutup tanpa pembuka", ">}]"},
		{"tutup dulu baru buka", "}{"},
		{"karakter bukan bracket", "<a>"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if ValidateBrackets(tt.in) {
				t.Errorf("ValidateBrackets(%q) = true, want false", tt.in)
			}
		})
	}
}

func TestValidateBrackets_Boundary(t *testing.T) {
	if ValidateBrackets("") {
		t.Error("empty string should be false (panjang minimum 1)")
	}
	tooLong := make([]byte, 4097)
	for i := range tooLong {
		tooLong[i] = '['
	}
	if ValidateBrackets(string(tooLong)) {
		t.Error("string > 4096 chars should be false")
	}
}