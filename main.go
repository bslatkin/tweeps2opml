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
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rodreegez/go-signin-with-twitter"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	socket, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	log.Printf("Listening on port %s\n", port)

	r := mux.NewRouter()
	r.HandleFunc("/", homepageHandler)
	r.HandleFunc("/list", listFriendsHandler)
	r.HandleFunc("/download", downloadHandler)
	r.HandleFunc("/authorize", signin.AuthorizeHandler)
	r.HandleFunc("/oauth_callback", signin.OauthCallbackHandler)

	server := &http.Server{Handler: r}
	if err := server.Serve(socket); err != nil {
		panic(err)
	}
}
