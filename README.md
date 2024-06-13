# Go LastFM Crawler
A command line tool to gather information and video links of artists from the lastfm website.

## Building locally
1. Install go locally
2. Clone https://github.com/steerpike/go-crawl.git
3. `cd crawl`
4. `go run . https://www.last.fm/music/The+Beatles`





# Music Crawler
Music Crawler is a command-line application written in Go that uses web scraping to gather information about music artists from a given URL. It uses the Colly library for web scraping and SQLite for data storage.

## Features
* Fetch artist details from a given URL
* Harvest a list of URLs to find artist details
* Find related artists to add to the seed list
* Store artist details in a SQLite database

### Commands
`fetch`: Fetch a URL for artist details and to find other related artists
`harvest`: Harvest the seed list for artist details and to find other related artists
Usage
To fetch artist details from a URL:
`go run main.go fetch <url>`
To harvest the seed list for artist details:
`go run main.go harvest`
