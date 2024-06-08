package util

import (
	"strings"
)

func GetId(url string) string {
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func GetNumPage(url string) int {
	pos := strings.Index(url, "trang-")
	if pos == -1 {
		return 1
	}

	pos += len("trang-")
	numPage := 0

	for i := pos; i < len(url); i++ {
		if url[i] >= '0' && url[i] <= '9' {
			numPage = numPage*10 + int(url[i]-'0')
		} else {
			break
		}

	}

	if numPage == 0 {
		return 1
	}
	return numPage
}
