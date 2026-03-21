package handlers

import "strings"

var logSanitizer = strings.NewReplacer("\n", "", "\r", "")

// sanitizeLog strips newline and carriage-return characters to prevent log injection.
func sanitizeLog(s string) string {
	return logSanitizer.Replace(s)
}
