// Package main implements Soal 1: case-insensitive string matching.
//
// Baca N lalu N string dari stdin. Cari SET PERTAMA string yang punya
// duplikat case-insensitive, lalu print semua index 1-based dipisah spasi.
// Kalau tidak ada duplikat sama sekali, print "false".
//
// Tanpa strings.ToLower/EqualFold/Compare, reflect.DeepEqual, atau helper
// search slice built-in. Case folding pakai bit manipulation ASCII.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	fmt.Println(solve(os.Stdin))
}

// solve membaca input dari r lalu mengembalikan jawaban sebagai string.
// Format: baris pertama N, lalu N baris string.
func solve(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	if !scanner.Scan() {
		return "false"
	}
	n, err := strconv.Atoi(trimCR(scanner.Text()))
	if err != nil || n <= 0 {
		return "false"
	}

	strs := make([]string, 0, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			break
		}
		strs = append(strs, trimCR(scanner.Text()))
	}

	for i := 0; i < len(strs); i++ {
		matches := []int{i + 1}
		for j := i + 1; j < len(strs); j++ {
			if equalIgnoreCase(strs[i], strs[j]) {
				matches = append(matches, j+1)
			}
		}
		if len(matches) > 1 {
			return joinInts(matches)
		}
	}

	return "false"
}

// equalIgnoreCase membandingkan dua string secara case-insensitive tanpa
// memakai strings.EqualFold/ToLower. Huruf A-Z dan a-z dibedakan oleh
// bit 0x20, jadi OR dengan 0x20 cukup untuk normalisasi ASCII.
func equalIgnoreCase(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca := a[i]
		cb := b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca |= 0x20
		}
		if cb >= 'A' && cb <= 'Z' {
			cb |= 0x20
		}
		if ca != cb {
			return false
		}
	}
	return true
}

// joinInts menggabungkan int ke string dipisah spasi tanpa strings.Join.
func joinInts(nums []int) string {
	out := ""
	for i := 0; i < len(nums); i++ {
		if i > 0 {
			out += " "
		}
		out += strconv.Itoa(nums[i])
	}
	return out
}

// trimCR buang trailing \r supaya input CRLF dari Windows tetap match.
func trimCR(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\r' {
		return s[:len(s)-1]
	}
	return s
}
