package util

import "time"

// JakartaLayout adalah format timestamp di JSON response (WIB).
const JakartaLayout = "2006-01-02 15:04:05"

// FormatJakarta mencetak t dalam zona Asia/Jakarta dengan JakartaLayout.
// Kalau zona gagal di-load (mis. tzdata tidak tersedia), fallback ke UTC.
func FormatJakarta(t time.Time) string {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return t.UTC().Format(JakartaLayout)
	}
	return t.In(loc).Format(JakartaLayout)
}