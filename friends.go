package main

import (
	"log"
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
</head>
<body>
    <div class="header">
        <h1>Tweeps 2 OPML</h1>
        <h2>Found {{. | len}} friends with URLs</h2>
    </div>
    <form action="/download" method="POST" target="_blank">
    {{range .}}
        <input name="{{.ScreenName}}" value="{{.URL}}" type="hidden" />
    {{end}}
    <input class="submit" type="submit" value="Download feeds" />
    </form>
</body>
</html>
`))
)

func listFriends(api *anaconda.TwitterApi) ([]anaconda.User, error) {
	result := make([]anaconda.User, 0, 1000)

	values := url.Values{}
	values.Set("cursor", "-1")
	values.Set("skip_status", "true")
	values.Set("include_user_entities", "false")
	values.Set("count", "200")

	// We're only allowed 15 requests every 15 minutes, so that's the most we'll even attempt to do. That means the most contacts you can export is 15 * 200 = 3000.
	for i := 0; i < 15; i++ {
		cursor, err := api.GetFriendsList(values)
		if err != nil {
			return []anaconda.User{}, err
		}
		if len(cursor.Users) == 0 {
			break
		}
		for _, user := range cursor.Users {
			if user.URL == "" {
				continue
			}
			result = append(result, user)
		}
		values.Set("cursor", cursor.Next_cursor_str)
	}

	return result, nil
}

func listFriendsHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := signin.GetUserInfo(r)
	if userInfo == nil {
		if err != nil {
			log.Printf("Could not get access token: %s", err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_SECRET"))
	api := anaconda.NewTwitterApi(userInfo.Token, userInfo.Secret)

	allFriends, err := listFriends(api)
	if err != nil {
		log.Printf("Could not list friends: %s", err)
		http.Error(w, "Could not list friends", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	listFriendsTemplate.Execute(w, allFriends)
}
