# Tweeps 2 OPML

*Get the RSS feeds of all of my Twitter friends*

This tool will find all of the RSS feeds that your Twitter friends provide and put them all into a single OPML file that you can then import into NewsBlur, The Old Reader, Feedly, Digg, etc.

## Developing

First, you need to create a `.env` file in your directory with the various Twitter authentication keys for your deployment. See `.env_example` for details.

To run the server:

```
go run *.go
```

## Future ideas

- Follow all rel="me" links and try to find feeds there too
- Sites like about.me don't include those links in some cases; special case them?
