package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/gorilla/mux"
	"github.com/rodreegez/go-signin-with-twitter"
)

var (
	notAuthenticatedTemplate = template.Must(template.New("").Parse(`
<html><body>
Generate an OPML file of my Twitter friends
<form action="/authorize" method="POST">
<input type="submit" value="Sign in to get started" />
<input type="hidden" name="continue" value="/generate" />
<input type="hidden" name="callback" value="/oauth_callback" />
</form>
</body></html>
`))
	friendsListTemplate = template.Must(template.New("").Parse(`
<html><body>
{{range .}}
	{{if .URL}}
		{{.ScreenName}} = {{.URL}}
		<br>
	{{end}}
{{else}}
No friends
{{end}}
</body></html>
`))
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	notAuthenticatedTemplate.Execute(w, nil)
}

func listAllFriends(api *anaconda.TwitterApi) ([]anaconda.User, error) {
	result := make([]anaconda.User, 0, 1000)

	values := url.Values{}
	values.Set("cursor", "-1")
	values.Set("skip_status", "true")
	values.Set("include_user_entities", "false")
	values.Set("count", "200")

	for {
		cursor, err := api.GetFriendsList(values)
		if err != nil {
			log.Printf("Could not list friends with arguments=%#v err=%s", values, err)
			return []anaconda.User{}, err
		}
		if len(cursor.Users) == 0 {
			break
		}
		result = append(result, cursor.Users...)
		values.Set("cursor", cursor.Next_cursor_str)
	}

	return result, nil
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
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

	allFriends, err := listAllFriends(api)
	if err != nil {
		http.Error(w, "Could not list friends", http.StatusInternalServerError)
		return
	}
	friendsListTemplate.Execute(w, allFriends)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homePageHandler)
	r.HandleFunc("/generate", generateHandler)
	r.HandleFunc("/authorize", signin.AuthorizeHandler)
	r.HandleFunc("/oauth_callback", signin.OauthCallbackHandler)
	server := &http.Server{Handler: r}
	port := os.Getenv("SERVER_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if nil != err {
		log.Fatalln(err)
	}
	log.Printf("Listening on port %s\n", port)
	if err := server.Serve(listener); nil != err {
		log.Fatalln(err)
	}
}
