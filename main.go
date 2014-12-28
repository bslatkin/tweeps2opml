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
	"path/filepath"

	"github.com/bslatkin/go-signin-with-twitter"
	"github.com/gorilla/context"
)

var (
	Globals = map[string]string{
		"AnalyticsId": os.Getenv("ANALYTICS_ID"),
	}
)

const (
	logsTarget = "/var/log/app_engine/custom_logs/tweeps2opml.log"
)

type Context struct {
	Globals interface{}
	Data    interface{}
}

func init() {
	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/list", listFriendsHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/authorize", signin.AuthorizeHandler)
	http.HandleFunc("/oauth_callback", signin.OauthCallbackHandler)

	// Configure logger to write to /var/log/app_engine/custom_logs if in the App Engine environment so it gets picked up in the admin console. See https://cloud.google.com/appengine/articles/logging
	if os.Getenv("APPENGINE_PRODUCTION") != "" {
		log.Printf("Attempting to write logs to a file")

		if err := os.MkdirAll(filepath.Dir(logsTarget), 0666); err != nil {
			log.Printf("Could not create logs directory: %s", err)
			return
		}

		writer, err := os.OpenFile(logsTarget, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Could not set log default output to file: %s", err)
			return
		}
		log.SetOutput(writer)

		log.Printf("Writing logs to file")
	}
}

func main() {
	port := os.Getenv("SERVER_PORT")
	log.Printf("Listening on port %s\n", port)
	// Wrap the DefaultServerMux to prevent memory leaks from Gorilla.
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), context.ClearHandler(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
