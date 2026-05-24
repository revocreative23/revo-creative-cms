package handlers

import (
	"os"
	"regexp"
)

var subdirSafe = regexp.MustCompile(`^[a-z0-9_-]+$`)

// sanitizeSubdir hanya izinkan huruf kecil, angka, underscore, dash.
// Mencegah path traversal seperti "../etc" atau "/abs/path".
func sanitizeSubdir(s string) string {
	if s == "" {
		return ""
	}
	if !subdirSafe.MatchString(s) {
		return ""
	}
	return s
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
