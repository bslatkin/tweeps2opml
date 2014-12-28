# Tweeps 2 OPML

*Get the RSS feeds of all of my Twitter friends*

This tool will find all of the RSS feeds that your Twitter friends provide and put them all into a single OPML file that you can then import into NewsBlur, The Old Reader, Feedly, Digg, etc.

## Developing

First, you need to create a `.env` file in your directory with the various Twitter authentication keys for your deployment. See `.env_example` for details.

To run the server locally:

```
go run *.go
```

To run in App Engine locally, follow the [getting started guide](https://cloud.google.com/appengine/docs/go/managed-vms/), then run:

```
gcloud preview app run app.yaml
```

If everything goes right, you should be able to hit <http://localhost:8080/> and see the app running.

To deploy to production, first set your project:

```
gcloud config set project <your-project-id>
```

Then deploy:

```
gcloud preview app deploy app.yaml
```

## Future ideas

- Follow all rel="me" links and try to find feeds there too
- Sites like about.me don't include those links in some cases; special case them?
- Do the same thing for Facebook, LinkedIn, Google+, etc
