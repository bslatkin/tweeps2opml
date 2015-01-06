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
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/ChimeraCoder/anaconda"
	"github.com/rodreegez/go-signin-with-twitter"
)

var (
	listFriendsTemplate = template.Must(template.New("friends").Parse(`
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Tweeps 2 OPML &raquo; Friends found</title>
    <style>
        body {
            font-family: Helvetica, Arial, sans-serif;
            max-width: 600px;
        }
        h1 {
            font-size: 32px;
        }
        h2 {
            font-size: 24px;
            font-weight: normal;
        }
        .submit {
            font-size: 24px;
            padding: 10px;
            background-color: #eee;
            border: 1px solid #ccc;
            border-radius: 4px;
        }
    </style>
    {{if .Globals.AnalyticsId}}
    <script>
      (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
      (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
      m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
      })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

      ga('create', '{{.Globals.AnalyticsId}}', 'auto');
      ga('send', 'pageview');
    </script>
    {{end}}
</head>
<body>
    <div class="header">
        <h1>Tweeps 2 OPML</h1>
        <h2>Found {{.Data | len}} friends with URLs</h2>
    </div>
    <div>
        <p>
            This download can take a while because we have to fetch many web pages and parse their HTML content to discover feeds. For 1000+ friends this may take up to 10 minutes. So please be patient!
        </p>
        <p>
            Once the download finishes, you should be able to import the resulting OPML file into your feed reader. Here's how to do it for <a href="http://www.onebigfluke.com/2014/12/how-to-import-opml-file-of-rss-feeds.html">NewsBlur</a>, <a href="http://blog.feedly.com/2013/07/03/the-fix-to-the-missing-feeds-issue-is-here/">Feedly</a>, <a href="http://theoldreader.com/pages/tour">The Old Reader</a>, and <a href="http://www.cnet.com/how-to/how-to-import-your-google-reader-data-to-digg-reader/">Digg Reader</a>.
        </p>
    </div>
    <form action="/download" method="POST" target="_blank">
    {{range .Data}}
        <input name="{{.ScreenName}}" value="{{.ProfileUrl}}" type="hidden" />
    {{end}}
        <input class="submit" type="submit" value="Download feeds" />
    </form>
</body>
</html>
`))
)

type Friend struct {
	ScreenName string
	ProfileUrl string
}

func fakeListFriends() ([]Friend, error) {
	return []Friend{
		// Both Atom and RSS
		Friend{"haxor", "http://onebigfluke.com"},
		// Relative feed URL
		Friend{"t", "http://tantek.com"},
		// A lot of feeds with the same content
		Friend{"adactio", "http://adactio.com"},
		// No feed title
		Friend{"polvi", "http://alex.polvi.net"},
		// Wikipedia link
		Friend{"evanpro", "https://en.wikipedia.org/wiki/Evan_Prodromou"},
		// Usually has EOF errors because of an old web server
		Friend{"markos", "http://www.dailykos.com"},
	}, nil
}

func listFriends(api *anaconda.TwitterApi) ([]Friend, error) {
	// NOTE: Use this for local development without hitting the Twitter API
	return fakeListFriends()

	result := make([]Friend, 0, 1000)

	values := url.Values{}
	values.Set("cursor", "-1")
	values.Set("skip_status", "true")
	values.Set("include_user_entities", "false")
	values.Set("count", "200")

	// We're only allowed 15 requests every 15 minutes, so that's the most we'll even attempt to do. That means the most contacts you can export is 15 * 200 = 3000.
	for i := 0; i < 15; i++ {
		cursor, err := api.GetFriendsList(values)
		if err != nil {
			return []Friend{}, err
		}
		if len(cursor.Users) == 0 {
			break
		}
		for _, user := range cursor.Users {
			if user.URL == "" {
				continue
			}
			result = append(result, Friend{user.ScreenName, user.URL})
		}
		values.Set("cursor", cursor.Next_cursor_str)
	}

	return result, nil
}

func listFriendsHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	userInfo, err := signin.GetUserInfo(r)
	if userInfo == nil {
		if err != nil {
			c.Log("Could not get access token: %s", err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_SECRET"))
	api := anaconda.NewTwitterApi(userInfo.Token, userInfo.Secret)

	allFriends, err := listFriends(api)
	if err != nil {
		c.Log("Could not list friends: %s", err)
		http.Error(w, "Could not list friends", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	listFriendsTemplate.Execute(w, Params{Globals: Globals, Data: allFriends})
}
