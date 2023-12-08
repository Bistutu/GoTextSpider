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

	var parseText func(*html.Node, []string) []string
	parseText = func(n *html.Node, texts []string) []string {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "img", "noscript":
				return texts
			}
		}
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if len(text) > 0 {
				texts = append(texts, text)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			texts = parseText(c, texts)
		}
		return texts
	}

	return parseText(doc, make([]string, 0)), nil
}
