// Package main implements Soal 3: validasi bracket tanpa regex.
//
// ValidateBrackets memastikan:
//   - Setiap pembuka (<, {, [) ditutup dengan pasangan yang benar
//   - Tidak ada penutup sebelum pembuka (misal "]<>")
//   - Tidak ada bracket yang "mengurung" bracket lain secara tidak penuh
//     (misal "<[>]" — saat ">" tiba, top stack adalah "[", mismatch)
//
// Algoritma: stack manual. Push pembuka; saat ketemu penutup, pop dan
// pastikan pasangannya benar. Di akhir string, stack harus kosong.
package main

import "fmt"

// ValidateBrackets mengembalikan true kalau string s punya bracket yang
// seimbang dan properly nested. Input hanya boleh mengandung <>{}[].
// Panjang valid: 1..4096 — di luar range langsung false.
func ValidateBrackets(s string) bool {
	if len(s) < 1 || len(s) > 4096 {
		return false
	}

	stack := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '<', '{', '[':
			stack = append(stack, c)
		case '>':
			if len(stack) == 0 || stack[len(stack)-1] != '<' {
				return false
			}
			stack = stack[:len(stack)-1]
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				return false
			}
			stack = stack[:len(stack)-1]
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return false
			}
			stack = stack[:len(stack)-1]
		default:
			return false
		}
	}
	return len(stack) == 0
}

func main() {
	samples := []string{
		"{{[<>[{{}}]]}}",
		"{<{[[{{[]<{{[{[]<>}]}}<>>}}]]}>}",
		"[{}<>]",
		"]",
		"][",
		"[>]",
		"[>",
		"[{}<[>]",
	}
	for i := 0; i < len(samples); i++ {
		fmt.Printf("%-40q => %v\n", samples[i], ValidateBrackets(samples[i]))
	}
}