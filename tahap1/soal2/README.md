# Soal 2 — Kembalian Kasir

`CalculateChange(total, paid)` mengembalikan kembalian + breakdown pecahan
(100.000..100). Kembalian dibulatkan ke bawah kelipatan Rp 100. Jika `paid < total`
return `(false, 0, 0, nil)` — pieces eksplisit `nil` (bukan map kosong) supaya
caller bisa membedakan "uang kurang" vs "tidak ada kembalian".

## Run

```bash
go run .
go test ./...
```

## Contoh

Total 700.649, bayar 800.000 -> kembalian 99.351 (bulat 99.300):
- 1 lembar 50.000, 2 lembar 20.000, 1 lembar 5.000, 2 lembar 2.000, 1 koin 200, 1 koin 100

Total 575.650, bayar 580.000 -> kembalian 4.350 (bulat 4.300):
- 2 lembar 2.000, 1 koin 200, 1 koin 100

Total 657.650, bayar 600.000 -> `false, kurang bayar`