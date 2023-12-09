package main

import (
	"fmt"

	"GoTextSpider/common"
)

func main() {
	url := "https://github.com/"
	fmt.Println(common.FetchAndParse(url))
	//fmt.Println(common.Sniffer(url))
}
