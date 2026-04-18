// Package main implements Soal 2: perhitungan kembalian kasir.
//
// CalculateChange menerima total belanja dan jumlah dibayar, lalu
// mengembalikan kembalian + breakdown pecahan. Pembulatan kembalian
// selalu ke bawah kelipatan Rp 100 (tidak ada pecahan < 100).
package main

import (
	"fmt"
	"strconv"
)

// denominations diurutkan dari besar ke kecil supaya algoritma greedy
// langsung menghasilkan jumlah lembar/koin minimum.
var denominations = []int64{
	100000, 50000, 20000, 10000, 5000,
	2000, 1000, 500, 200, 100,
}

// CalculateChange menghitung kembalian dari (total, paid).
//
// Return:
//   - ok:       false kalau paid < total (uang kurang)
//   - change:   kembalian asli (paid - total)
//   - rounded:  kembalian setelah dibulatkan ke bawah Rp 100
//   - pieces:   map[pecahan]jumlah_lembar_atau_koin (hanya yang > 0)
func CalculateChange(total, paid int64) (bool, int64, int64, map[int64]int64) {
	if paid < total {
		// Semua return value zero-value eksplisit saat uang kurang:
		// bool=false, int64=0, map=nil (bukan map kosong) supaya caller
		// bisa membedakan "tidak ada kembalian" vs "input invalid".
		var pieces map[int64]int64
		return false, 0, 0, pieces
	}

	change := paid - total
	rounded := (change / 100) * 100

	pieces := make(map[int64]int64)
	remaining := rounded
	for i := 0; i < len(denominations); i++ {
		d := denominations[i]
		if remaining >= d {
			count := remaining / d
			pieces[d] = count
			remaining -= count * d
		}
	}

	return true, change, rounded, pieces
}

// formatPieces merangkai breakdown pecahan jadi list baris manusia-baca.
// Mengikuti urutan denominations (besar ke kecil) supaya output
// deterministik — map iteration order di Go acak, jadi kita tidak iterasi map.
func formatPieces(pieces map[int64]int64) []string {
	out := []string{}
	for i := 0; i < len(denominations); i++ {
		d := denominations[i]
		n, ok := pieces[d]
		if !ok || n == 0 {
			continue
		}
		unit := "lembar"
		if d < 1000 {
			unit = "koin"
		}
		out = append(out, strconv.FormatInt(n, 10)+" "+unit+" "+formatRupiah(d))
	}
	return out
}

// formatRupiah mencetak angka dengan pemisah ribuan '.'.
// Contoh: 50000 -> "50.000", 200 -> "200".
func formatRupiah(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}
	out := ""
	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		if count == 3 {
			out = "." + out
			count = 0
		}
		out = string(s[i]) + out
		count++
	}
	return out
}

func main() {
	cases := []struct {
		total, paid int64
	}{
		{700649, 800000},
		{575650, 580000},
		{657650, 600000},
	}

	for i := 0; i < len(cases); i++ {
		c := cases[i]
		fmt.Printf("Total belanja seorang customer: %d, Pembeli membayar: %d\n", c.total, c.paid)
		ok, change, rounded, pieces := CalculateChange(c.total, c.paid)
		if !ok {
			fmt.Println("  false (uang kurang)")
			fmt.Println()
			continue
		}
		fmt.Printf("  Kembalian yang harus diberikan kasir: %d (dibulatkan: %d)\n", change, rounded)
		lines := formatPieces(pieces)
		for j := 0; j < len(lines); j++ {
			fmt.Println("  - " + lines[j])
		}
		fmt.Println()
	}
}
