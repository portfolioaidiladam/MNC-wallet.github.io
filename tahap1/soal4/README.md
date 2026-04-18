# Soal 4 — Validasi Cuti Pribadi

`CanTakeLeave(cutiBersama, joinDate, planDate, duration)` mengembalikan
`(bool, alasan)`. Aturan: cuti kantor 14 hari/tahun, cuti pribadi = 14 -
cutiBersama, karyawan baru tidak boleh cuti 180 hari pertama, max 3 hari
berturutan, kuota tahun pertama proporsional terhadap sisa hari kerja.

## Run

```bash
go run .
go test ./...
```

## Contoh (PDF)

- Join 2021-05-01, plan 2021-07-05, dur 1 -> `false`, "belum 180 hari..."
- Join 2021-05-01, plan 2021-11-05, dur 3 -> `false`, "hanya boleh mengambil 1 hari cuti"
- Join 2021-01-05, plan 2021-12-18, dur 1 -> `true`
- Join 2021-01-05, plan 2021-12-18, dur 3 -> `true`

## Asumsi Implementasi

PDF hanya menspesifikasikan perhitungan kuota untuk **tahun pertama** (prorated
dari `(joinDate + 180) s/d 31 Desember`). Implementasi ini mengasumsikan **tahun
kedua dan seterusnya dapat kuota penuh** `14 - cutiBersama` hari. Kalau
interpretasi berbeda (mis. reset kuota di anniversary date, atau prorata lagi
tiap tahun), logic di `quotaForYear` — khususnya cabang `planDate.Year() !=
joinDate.Year()` — perlu disesuaikan.