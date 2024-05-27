package main

import (
	"fmt"
	"novel_crawler/internal"
)

func main() {

	truyenFull := &internal.TruyenFull{}
	listCategories := truyenFull.GetCategories()

	for _, v := range listCategories {
		fmt.Println(v)
	}

}
