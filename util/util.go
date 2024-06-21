package util

import (
	"strconv"
	"strings"
)

func Max(num1, num2 int) int {
	if num1 > num2 {
		return num1
	}
	return num2
}

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

func FindPrevAndNextChapters(current string, chapterNew string, chapterLast string) (string, string) {

	temps := strings.Split(current, "chuong-")
	temp := temps[len(temps)-1]

	num, _ := strconv.Atoi(temp)

	if current == chapterLast || num == 1 {
		return current, "chuong-" + strconv.Itoa(num+1)
	}

	if current == chapterNew {
		return "chuong-" + strconv.Itoa(num-1), current
	}
	prev := strconv.Itoa(num - 1)
	next := strconv.Itoa(num + 1)
	prev = "chuong-" + prev
	next = "chuong-" + next
	return prev, next
}
