// Package main implements Soal 4: validasi pengajuan cuti pribadi karyawan.
//
// Aturan:
//   - Cuti kantor 14 hari/tahun, cuti pribadi = 14 - cuti_bersama.
//   - Karyawan baru: tidak boleh cuti pribadi selama 180 hari pertama.
//   - Kuota tahun pertama = floor(hari_tersisa / 365 * cuti_pribadi),
//     hari_tersisa = dari (joinDate + 180 hari) s/d 31 Desember tahun itu.
//   - Cuti pribadi max 3 hari berturutan per pengambilan.
package main

import (
	"fmt"
	"time"
)

// CanTakeLeave memvalidasi apakah pengajuan cuti pribadi bisa disetujui.
//
// Urutan cek:
//  1. cutiBersama wajib di rentang [0, 14].
//  2. duration wajib > 0.
//  3. planDate wajib >= joinDate + 180 hari.
//  4. duration wajib <= 3 (maksimal 3 hari berturutan).
//  5. duration wajib <= kuota cuti tahun tersebut.
//
// Return (true, "") kalau semua cek lolos. Kalau gagal, return
// (false, alasan).
func CanTakeLeave(cutiBersama int, joinDate, planDate time.Time, duration int) (bool, string) {
	if cutiBersama < 0 || cutiBersama > 14 {
		return false, "Jumlah cuti bersama tidak valid"
	}
	if duration <= 0 {
		return false, "Durasi cuti tidak valid"
	}

	joinDate = truncateToDay(joinDate)
	planDate = truncateToDay(planDate)
	activeDate := joinDate.AddDate(0, 0, 180)

	if planDate.Before(activeDate) {
		return false, "Karena belum 180 hari sejak tanggal join karyawan"
	}

	if duration > 3 {
		return false, "Karena cuti pribadi maksimal 3 hari berturutan"
	}

	personalLeave := 14 - cutiBersama
	quota := quotaForYear(joinDate, planDate, activeDate, personalLeave)

	if duration > quota {
		return false, fmt.Sprintf("Karena hanya boleh mengambil %d hari cuti", quota)
	}

	return true, ""
}

// quotaForYear menghitung kuota cuti pribadi untuk tahun planDate.
//
// Kalau planDate masih di tahun kalender yang sama dengan joinDate,
// pakai formula proporsional: floor(daysRemaining / 365 * personalLeave).
// Tahun-tahun berikutnya karyawan sudah dapat kuota penuh.
func quotaForYear(joinDate, planDate, activeDate time.Time, personalLeave int) int {
	if planDate.Year() != joinDate.Year() {
		return personalLeave
	}

	dec31 := time.Date(joinDate.Year(), 12, 31, 0, 0, 0, 0, joinDate.Location())
	daysRemaining := int(dec31.Sub(activeDate).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}
	return (daysRemaining * personalLeave) / 365
}

// truncateToDay membuang komponen jam/menit/detik biar perbandingan
// tanggal tidak terpengaruh clock time. Zone dipertahankan.
func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func main() {
	parse := func(s string) time.Time {
		t, _ := time.Parse("2006-01-02", s)
		return t
	}
	cases := []struct {
		cutiBersama    int
		join, plan     string
		duration       int
	}{
		{7, "2021-05-01", "2021-07-05", 1},
		{7, "2021-05-01", "2021-11-05", 3},
		{7, "2021-01-05", "2021-12-18", 1},
		{7, "2021-01-05", "2021-12-18", 3},
	}
	for i := 0; i < len(cases); i++ {
		c := cases[i]
		ok, reason := CanTakeLeave(c.cutiBersama, parse(c.join), parse(c.plan), c.duration)
		fmt.Printf("join=%s plan=%s dur=%d => ok=%v reason=%q\n",
			c.join, c.plan, c.duration, ok, reason)
	}
}