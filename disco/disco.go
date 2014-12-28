package disco

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Feed struct {
	Title string
	URL   *url.URL
}

func traverse(rootUrl *url.URL, node *html.Node) []Feed {
	var result []Feed

	if node.Type == html.ElementNode && node.Data == "link" {
		var rel, linkType, title, href string
		for _, attr := range node.Attr {
			trimmed := Strings.TrimSpace(attr.Val)
			switch attr.Key {
			case "rel":
				rel = trimmed
			case "type":
				linkType = trimmed
			case "title":
				title = trimmed
			case "href":
				href = trimmed
			}
		}
		if rel == "alternate" && href != "" && (linkType == "application/atom+xml" || linkType == "application/rss+xml") {
			if parsed, err := url.Parse(href); err == nil {
				result = append(result, Feed{
					Title: title,
					URL:   rootUrl.ResolveReference(parsed),
				})
			}
		}
		// TODO: Consider traversing rel="me" links
	}

	if node.FirstChild != nil {
		result = append(result, traverse(rootUrl, node.FirstChild)...)
	}

	if node.PrevSibling == nil {
		for sibling := node.NextSibling; sibling != nil; sibling = sibling.NextSibling {
			result = append(result, traverse(rootUrl, sibling)...)
		}
	}

	return result
}

func Discover(url *url.URL) (feeds []Feed, err error) {
	var resp *http.Response
	for i := 0; i < 3; i++ {
		// TODO: Handle EOF errors here by using req.Header.Add("Accept-Encoding", "identity")
		resp, err = http.Get(url.String())
		if err == nil {
			break
		}
		time.Sleep(time.Duration(3) * time.Second)
	}
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	feeds = traverse(url, doc)
	return
}
