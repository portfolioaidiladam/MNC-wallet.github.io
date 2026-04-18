package main

import (
	"testing"
	"time"
)

func mustParseDate(t *testing.T, s string) time.Time {
	t.Helper()
	tm, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatalf("parse %q: %v", s, err)
	}
	return tm
}

func TestCanTakeLeave_PDFCases(t *testing.T) {
	tests := []struct {
		name        string
		cutiBersama int
		joinDate    string
		planDate    string
		duration    int
		wantOk      bool
		wantReason  string
	}{
		{
			name:        "PDF 1: rencana sebelum 180 hari",
			cutiBersama: 7,
			joinDate:    "2021-05-01",
			planDate:    "2021-07-05",
			duration:    1,
			wantOk:      false,
			wantReason:  "Karena belum 180 hari sejak tanggal join karyawan",
		},
		{
			name:        "PDF 2: sudah 180 hari tapi minta 3, kuota cuma 1",
			cutiBersama: 7,
			joinDate:    "2021-05-01",
			planDate:    "2021-11-05",
			duration:    3,
			wantOk:      false,
			wantReason:  "Karena hanya boleh mengambil 1 hari cuti",
		},
		{
			name:        "PDF 3: join awal tahun, kuota cukup, minta 1 hari",
			cutiBersama: 7,
			joinDate:    "2021-01-05",
			planDate:    "2021-12-18",
			duration:    1,
			wantOk:      true,
			wantReason:  "",
		},
		{
			name:        "PDF 4: join awal tahun, kuota cukup, minta 3 hari",
			cutiBersama: 7,
			joinDate:    "2021-01-05",
			planDate:    "2021-12-18",
			duration:    3,
			wantOk:      true,
			wantReason:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, reason := CanTakeLeave(
				tt.cutiBersama,
				mustParseDate(t, tt.joinDate),
				mustParseDate(t, tt.planDate),
				tt.duration,
			)
			if ok != tt.wantOk {
				t.Errorf("ok = %v, want %v", ok, tt.wantOk)
			}
			if reason != tt.wantReason {
				t.Errorf("reason = %q, want %q", reason, tt.wantReason)
			}
		})
	}
}

func TestCanTakeLeave_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		cutiBersama int
		joinDate    string
		planDate    string
		duration    int
		wantOk      bool
		wantReason  string
	}{
		{
			name:        "durasi > 3 hari ditolak",
			cutiBersama: 7,
			joinDate:    "2021-01-05",
			planDate:    "2021-12-18",
			duration:    4,
			wantOk:      false,
			wantReason:  "Karena cuti pribadi maksimal 3 hari berturutan",
		},
		{
			name:        "planDate tepat di hari ke-180 (boundary)",
			cutiBersama: 7,
			joinDate:    "2021-05-01",
			planDate:    "2021-10-28",
			duration:    1,
			wantOk:      true,
			wantReason:  "",
		},
		{
			name:        "planDate sehari sebelum hari ke-180",
			cutiBersama: 7,
			joinDate:    "2021-05-01",
			planDate:    "2021-10-27",
			duration:    1,
			wantOk:      false,
			wantReason:  "Karena belum 180 hari sejak tanggal join karyawan",
		},
		{
			name:        "join 2021-12-31, plan 2022-06-28 (H-1 hari ke-180)",
			cutiBersama: 5,
			joinDate:    "2021-12-31",
			planDate:    "2022-06-28",
			duration:    1,
			wantOk:      false,
			wantReason:  "Karena belum 180 hari sejak tanggal join karyawan",
		},
		{
			name:        "tahun kedua setelah 180 hari - kuota penuh",
			cutiBersama: 7,
			joinDate:    "2021-12-31",
			planDate:    "2022-07-01",
			duration:    3,
			wantOk:      true,
			wantReason:  "",
		},
		{
			name:        "cuti bersama 0, personal leave 14, minta 3 di akhir tahun join",
			cutiBersama: 0,
			joinDate:    "2021-01-05",
			planDate:    "2021-12-18",
			duration:    3,
			wantOk:      true,
			wantReason:  "",
		},
		{
			name:        "cuti bersama 14, personal leave 0, selalu ditolak kuota",
			cutiBersama: 14,
			joinDate:    "2021-01-05",
			planDate:    "2021-12-18",
			duration:    1,
			wantOk:      false,
			wantReason:  "Karena hanya boleh mengambil 0 hari cuti",
		},
		{
			name:        "join di tahun kabisat, hitung tetap benar",
			cutiBersama: 7,
			joinDate:    "2020-01-05",
			planDate:    "2020-12-18",
			duration:    3,
			wantOk:      true,
			wantReason:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, reason := CanTakeLeave(
				tt.cutiBersama,
				mustParseDate(t, tt.joinDate),
				mustParseDate(t, tt.planDate),
				tt.duration,
			)
			if ok != tt.wantOk {
				t.Errorf("ok = %v, want %v", ok, tt.wantOk)
			}
			if reason != tt.wantReason {
				t.Errorf("reason = %q, want %q", reason, tt.wantReason)
			}
		})
	}
}

func TestQuotaForYear(t *testing.T) {
	// Verifikasi spec PDF: 1 Mei 2021 join, kuota = floor(64/365 * 7) = 1
	join := mustParseDate(t, "2021-05-01")
	plan := mustParseDate(t, "2021-11-05")
	active := join.AddDate(0, 0, 180)
	got := quotaForYear(join, plan, active, 7)
	if got != 1 {
		t.Errorf("quotaForYear(join=2021-05-01, personalLeave=7) = %d, want 1", got)
	}

	// Verifikasi: 5 Januari 2021 join -> active = 4 Juli 2021
	// days = Dec 31 - Jul 4 = 180. floor(180/365*7) = 3
	join2 := mustParseDate(t, "2021-01-05")
	plan2 := mustParseDate(t, "2021-12-18")
	active2 := join2.AddDate(0, 0, 180)
	got2 := quotaForYear(join2, plan2, active2, 7)
	if got2 != 3 {
		t.Errorf("quotaForYear(join=2021-01-05, personalLeave=7) = %d, want 3", got2)
	}
}