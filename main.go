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
	"net/http"
	"os"

	"github.com/bslatkin/go-signin-with-twitter"
	"github.com/gorilla/context"
)

func init() {
	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/list", listFriendsHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/authorize", signin.AuthorizeHandler)
	http.HandleFunc("/oauth_callback", signin.OauthCallbackHandler)

	// TODO: Configure logger to write to /var/log/app_engine/custom_logs if in the App Engine environment. See https://cloud.google.com/appengine/articles/logging
}

func main() {
	port := os.Getenv("SERVER_PORT")
	log.Printf("Listening on port %s\n", port)
	// Wrap the DefaultServerMux to prevent memory leaks from Gorilla.
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), context.ClearHandler(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
