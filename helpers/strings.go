package helpers

import "strings"

func SplitString(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	var b strings.Builder
	b.Grow(chunkSize)
	l := 0
	for _, r := range s {
		b.WriteRune(r)
		l++
		if l == chunkSize {
			chunks = append(chunks, b.String())
			l = 0
			b.Reset()
			b.Grow(chunkSize)
		}
	}
	if l > 0 {
		chunks = append(chunks, b.String())
	}
	return chunks
}

func ContainsStringInStringArray(arr []string, key string) bool {
	for _, str := range arr {
		if str == key {
			return true
		}
	}
	return false
}
