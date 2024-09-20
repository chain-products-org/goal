package mathx

import "strings"

func Trim0(s string) string {
	idx := strings.Index(s, ".")
	if idx < 0 {
		return s
	}
	prefix := s[:idx]
	suffix := s[idx+1:]
	suffix = strings.TrimRight(suffix, "0")
	if suffix == "" {
		return prefix
	}
	return prefix + "." + suffix
}
