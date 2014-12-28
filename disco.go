/*
 * Copyright 2014 Brett Slatkin
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Feed struct {
	Title string
	Url   string
}

func traverse(rootUrl *url.URL, node *html.Node) []Feed {
	var result []Feed

	if node.Type == html.ElementNode && node.Data == "link" {
		var rel, linkType, title, href string
		for _, attr := range node.Attr {
			trimmed := strings.TrimSpace(attr.Val)
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
					Url:   rootUrl.ResolveReference(parsed).String(),
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
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
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

type feedLengthSorter []Feed

func (f feedLengthSorter) Len() int {
	return len(f)
}

func (f feedLengthSorter) Less(i, j int) bool {
	return len(f[i].Url) < len(f[j].Url)
}

func (f feedLengthSorter) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func GetPrimaryFeed(feeds []Feed) Feed {
	comments := map[string]Feed{}
	content := map[string]Feed{}

	// Dedupe all the URLs.
	for _, feed := range feeds {
		if strings.Contains(feed.Url, "wikipedia.org") {
			// People like to link to wikipedia, but the only Atom feed they provide is *all* changes to all pages, not just the one linked to.
			continue
		} else if strings.Contains(feed.Url, "comments") {
			comments[feed.Url] = feed
		} else {
			content[feed.Url] = feed
		}
	}

	// Allow the comment feeds to be considered if there are no other feeds.
	if len(content) == 0 {
		for feedUrl, feed := range comments {
			content[feedUrl] = feed
		}
	}

	// No feeds found to consider.
	if len(content) == 0 {
		return Feed{}
	}

	filtered := make([]Feed, 0, len(content))
	for _, feed := range content {
		filtered = append(filtered, feed)
	}

	sort.Sort(feedLengthSorter(filtered))

	// Return the shortest URL, which is usually the primary feed, often Atom.
	return filtered[0]
}
