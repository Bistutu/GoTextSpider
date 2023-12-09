package common

import (
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"GoTextSpider/log"
)

// FetchAndParse 读取单个网页并解析文本
func FetchAndParse(link string) ([]string, error) {
	resp, err := http.Get(link)
	if err != nil {
		log.Errorf("fail to get link: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Errorf("fail to parse html: %v", err)
		return nil, err
	}
	rs := make([]string, 0)

	var parseText func(*html.Node)
	parseText = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "img", "noscript":
				return
			}
		}
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if len(text) > 0 {
				rs = append(rs, text)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseText(c)
		}
	}
	parseText(doc)

	return rs, nil
}

// Sniffer 读取单个网页并解析文本
func Sniffer(link string) error {
	queue := make([]string, 0)

	// 解析初始链接，获取主机名
	initialURL, _ := url.Parse(link)
	initialHost := initialURL.Host

	queue = append(queue, link)

	for len(queue) > 0 {
		link := queue[0]
		queue = queue[1:]

		resp, err := http.Get(link)
		if err != nil {
			continue
		}
		doc, err := html.Parse(resp.Body)
		if err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		var parseText func(string, *html.Node)
		parseText = func(link string, n *html.Node) {
			if n.Type == html.ElementNode {
				switch n.Data {
				case "script", "style", "img", "noscript":
					return
				case "a":
					for _, a := range n.Attr {
						if a.Key == "href" {
							href, err := url.Parse(a.Val)
							if err != nil {
								continue
							}
							// 如果主机名匹配，将链接加入队列
							if href.Host == initialHost {
								queue = append(queue, a.Val)
							}
						}
					}
				}
			}
			if n.Type == html.TextNode {
				text := strings.TrimSpace(n.Data)
				if len(text) > 0 {
					// TODO 处理逻辑
					// link : text
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				parseText(link, c)
			}
		}
		parseText(link, doc)
	}

	return nil
}
