package util

import (
	"strings"
)

func GetId(url string) string {
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func GetNumPage(url string, prefixes ...string) int {
	for _, prefix := range prefixes {
		pos := strings.Index(url, prefix)
		if pos == -1 {
			continue
		}
		pos += len(prefix)
		numPage := 0
		for i := pos; i < len(url); i++ {
			if url[i] >= '0' && url[i] <= '9' {
				numPage = numPage*10 + int(url[i]-'0')
			} else {
				break
			}
		}
		if numPage == 0 {
			continue
		}
		return numPage
	}
	return 1
}
