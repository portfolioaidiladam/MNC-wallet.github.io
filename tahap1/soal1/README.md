# Soal 1 — Case-Insensitive String Matching

Baca `N` lalu `N` string dari stdin. Temukan SET PERTAMA string yang punya duplikat
secara case-insensitive dan cetak semua index 1-based dipisah spasi. Jika tidak ada
duplikat, cetak `false`. Tidak memakai helper array bawaan; case folding pakai bit
manipulation ASCII (`|= 0x20`).

## Run

```bash
go run .
go test ./...
```

## Contoh

Input:
```
4
abcd
acbd
aaab
acbd
```
Output: `2 4`

Input:
```
11
Satu
Sate
Tujuh
Tusuk
Tujuh
Sate
Bonus
Tiga
Puluh
Tujuh
Tusuk
```
Output: `3 5 10`

Input: 5 string unik -> `false`