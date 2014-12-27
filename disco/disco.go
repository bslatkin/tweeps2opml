package disco

import (
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"time"
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
			switch attr.Key {
			case "rel":
				rel = attr.Val
			case "type":
				linkType = attr.Val
			case "title":
				title = attr.Val
			case "href":
				href = attr.Val
			}
		}
		if rel == "alternate" && href != "" && (linkType == "application/atom+xml" || linkType == "application/rss+xml") {
			// TODO: Trim the href
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
	if resp.Body != nil {
		// If the body is present then we'll read it. This gets around when the server closes the connection before we have a chance to read, which happens with servers that don't support HTTP keep-alive.
		err = nil
		defer resp.Body.Close()
	} else if err != nil {
		// The body is missing and there's an error. This is probably a DNS lookup failure.
		return
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	feeds = traverse(url, doc)
	return
}
