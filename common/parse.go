package common

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"GoTextSpider/log"
)

func FetchAndParse(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("fail to get url: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Errorf("fail to parse html: %v", err)
		return nil, err
	}

	// 保存所有文本
	texts := make([]string, 0)

	var parseText func(*html.Node) []string
	parseText = func(n *html.Node) []string {
		// 如果是文本节点，且不是 script、style、img、noscript 则获取文本
		if n.Type == html.TextNode &&
			n.Parent.Data != "script" &&
			n.Parent.Data != "style" &&
			n.Parent.Data != "img" &&
			n.Parent.Data != "noscript" {
			text := strings.TrimSpace(n.Data)
			if len(text) > 0 {
				texts = append(texts, text)
			}
		}
		// 遍历所有子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseText(c)
		}
		return texts
	}
	return parseText(doc), nil
}
