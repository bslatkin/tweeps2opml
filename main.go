package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rodreegez/go-signin-with-twitter"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homepageHandler)
	r.HandleFunc("/list", listFriendsHandler)
	r.HandleFunc("/download", downloadHandler)
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
