package main

import (
	"reflect"
	"testing"
)

func TestCalculateChange(t *testing.T) {
	tests := []struct {
		name        string
		total, paid int64
		wantOk      bool
		wantChange  int64
		wantRounded int64
		wantPieces  map[int64]int64
	}{
		{
			name:        "PDF case 1: 700.649 dibayar 800.000",
			total:       700649,
			paid:        800000,
			wantOk:      true,
			wantChange:  99351,
			wantRounded: 99300,
			wantPieces: map[int64]int64{
				50000: 1,
				20000: 2,
				5000:  1,
				2000:  2,
				200:   1,
				100:   1,
			},
		},
		{
			name:        "PDF case 2: 575.650 dibayar 580.000",
			total:       575650,
			paid:        580000,
			wantOk:      true,
			wantChange:  4350,
			wantRounded: 4300,
			wantPieces: map[int64]int64{
				2000: 2,
				200:  1,
				100:  1,
			},
		},
		{
			name:        "PDF case 3: uang kurang",
			total:       657650,
			paid:        600000,
			wantOk:      false,
			wantChange:  0,
			wantRounded: 0,
			wantPieces:  nil,
		},
		{
			name:        "pas bayar, tidak ada kembalian",
			total:       50000,
			paid:        50000,
			wantOk:      true,
			wantChange:  0,
			wantRounded: 0,
			wantPieces:  map[int64]int64{},
		},
		{
			name:        "kembalian 100 tepat",
			total:       900,
			paid:        1000,
			wantOk:      true,
			wantChange:  100,
			wantRounded: 100,
			wantPieces:  map[int64]int64{100: 1},
		},
		{
			name:        "kembalian < 100 dibulatkan jadi 0",
			total:       950,
			paid:        1000,
			wantOk:      true,
			wantChange:  50,
			wantRounded: 0,
			wantPieces:  map[int64]int64{},
		},
		{
			name:        "pecahan 100.000 terpakai",
			total:       100000,
			paid:        350000,
			wantOk:      true,
			wantChange:  250000,
			wantRounded: 250000,
			wantPieces: map[int64]int64{
				100000: 2,
				50000:  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, change, rounded, pieces := CalculateChange(tt.total, tt.paid)
			if ok != tt.wantOk {
				t.Errorf("ok = %v, want %v", ok, tt.wantOk)
			}
			if change != tt.wantChange {
				t.Errorf("change = %d, want %d", change, tt.wantChange)
			}
			if rounded != tt.wantRounded {
				t.Errorf("rounded = %d, want %d", rounded, tt.wantRounded)
			}
			if !reflect.DeepEqual(pieces, tt.wantPieces) {
				t.Errorf("pieces = %v, want %v", pieces, tt.wantPieces)
			}
		})
	}
}

// TestCalculateChange_UnderpaidReturnsNil memastikan kontrak eksplisit:
// saat paid < total, pieces HARUS nil (bukan map kosong). Ini memperkuat
// pembedaan "input invalid" vs "tidak ada kembalian" (change 0 -> map{}).
func TestCalculateChange_UnderpaidReturnsNil(t *testing.T) {
	ok, change, rounded, pieces := CalculateChange(657650, 600000)
	if ok {
		t.Error("ok = true, want false for underpaid")
	}
	if change != 0 || rounded != 0 {
		t.Errorf("change=%d rounded=%d, want both 0", change, rounded)
	}
	if pieces != nil {
		t.Errorf("pieces = %v, want nil (not empty map)", pieces)
	}
}

func TestFormatPieces(t *testing.T) {
	pieces := map[int64]int64{
		50000: 1,
		20000: 2,
		5000:  1,
		2000:  2,
		200:   1,
		100:   1,
	}
	got := formatPieces(pieces)
	want := []string{
		"1 lembar 50.000",
		"2 lembar 20.000",
		"1 lembar 5.000",
		"2 lembar 2.000",
		"1 koin 200",
		"1 koin 100",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("formatPieces() =\n%v\nwant\n%v", got, want)
	}
}

func TestFormatRupiah(t *testing.T) {
	tests := []struct {
		in   int64
		want string
	}{
		{100, "100"},
		{200, "200"},
		{500, "500"},
		{1000, "1.000"},
		{5000, "5.000"},
		{50000, "50.000"},
		{100000, "100.000"},
	}
	for _, tt := range tests {
		got := formatRupiah(tt.in)
		if got != tt.want {
			t.Errorf("formatRupiah(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}