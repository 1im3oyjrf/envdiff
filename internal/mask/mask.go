package mask

import "strings"

const defaultMask = "****"

// SecretKeyPatterns contains substrings that indicate a key holds a secret value.
var SecretKeyPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
	"PWD",
}

// IsSecret returns true if the key name suggests it contains a sensitive value.
func IsSecret(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range SecretKeyPatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// MaskValue returns the masked representation of a secret value.
// If the value is empty, it returns an empty string.
func MaskValue(value string) string {
	if value == "" {
		return ""
	}
	return defaultMask
}

// ApplyMask returns the value as-is if masking is disabled or the key is not
// a secret; otherwise it returns the masked value.
func ApplyMask(key, value string, enabled bool) string {
	if !enabled || !IsSecret(key) {
		return value
	}
	return MaskValue(value)
}
