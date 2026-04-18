# Soal 3 — Validasi Bracket

`ValidateBrackets(s)` mengecek string yang hanya berisi `<>{}[]` seimbang dan
properly nested. Setiap pembuka harus ditutup dengan pasangan yang cocok, tidak
boleh ada penutup sebelum pembuka, dan tidak boleh ada bracket yang mengurung
bracket lain secara tidak penuh. Panjang valid 1..4096. Implementasi pakai stack
manual (tanpa regex).

## Run

```bash
go run .
go test ./...
```

## Contoh

TRUE:
- `{{[<>[{{}}]]}}`
- `{<{[[{{[]<{{[{[]<>}]}}<>>}}]]}>}`
- `[{}<>]`

FALSE:
- `]`
- `][`
- `[>]`
- `[{}<[>]`
- `<[>]` (bracket silang)