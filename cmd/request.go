/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crawl/models"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

func UrlExists(path string) bool {
	var exists bool
	db, err := sql.Open("sqlite3", "music.db")
	if err != nil {
		log.Println("Something went wrong opening database:", err)
	}
	query := `SELECT 1 FROM Artists WHERE Path = ?`
	err = db.QueryRow(query, path).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("querying for url existence: %v", err)
	}
	return exists
}

func CrawlURLs(url string, crawlLimit int, fromId int64) {
	c := colly.NewCollector()
	pathSet := make(map[string]bool)
	artist := models.Artist{}
	tags := []string{}
	similar := []string{}
	videos := make(map[string]string)
	c.OnResponse(func(r *colly.Response) {
		artist.Response = r.StatusCode
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		foundPath := e.Attr("href")
		// match, _ := regexp.MatchString("^/music/[^/]+(/(\\+tracks|\\+similar))?/?$", foundURL)
		if strings.Contains(foundPath, "+free-music-downloads") {
			return
		}
		match, _ := regexp.MatchString("^/music/[^/?#]+/?$", foundPath)
		if match {
			pathSet[foundPath] = true
		}
	})

	// Scrape tags
	c.OnHTML(".tag", func(e *colly.HTMLElement) {
		tag := e.Text
		tags = append(tags, tag)
	})

	// Scrape similar artists
	c.OnHTML("h3.artist-similar-artists-sidebar-item-name a", func(e *colly.HTMLElement) {
		a := e.Text
		similar = append(similar, a)
	})

	// Scrape canonical name
	c.OnHTML("link[rel=canonical]", func(e *colly.HTMLElement) {
		canonicalName := e.Attr("href")
		artist.Url = canonicalName
	})

	// Scrape video names and urls
	c.OnHTML("td[class=chartlist-play] > a", func(e *colly.HTMLElement) {
		videoName := e.Attr("data-track-name")
		videoUrl := e.Attr("href")
		path := e.Attr("data-artist-url")
		artist.Path = path
		videos[videoUrl] = videoName
	})

	// Scrape artists name
	c.OnHTML("#tonefuze-mobile", func(e *colly.HTMLElement) {
		artist.Name = e.Attr("data-tonefuze-artist")
	})

	c.OnScraped(func(r *colly.Response) {
		artist.Tags = tags
		artist.Videos = videos
		artist.Similar = similar
		artist.FromID = fromId
		if strings.Contains(artist.Path, "music") {
			artist.Save()
		}

		for path := range pathSet {
			if crawlLimit > 0 {
				crawlLimit--
				if !UrlExists(path) {
					url := "https://www.last.fm" + path
					CrawlURLs(url, crawlLimit, artist.ID)
					fmt.Println("Url found and crawled:", url)
				}
			}
		}
	})

	c.Visit(url)
}

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl a URL and find other URLs to crawl",
	Long:  `Crawl a URL using the Colly library and find other URLs on that page to crawl`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		crawlLimit, _ := cmd.Flags().GetInt("limit")
		CrawlURLs(url, crawlLimit, 0)
	},
}

func init() {
	crawlCmd.Flags().IntP("limit", "l", 5, "Limit of how many URLs to crawl")
	rootCmd.AddCommand(crawlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// requestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// requestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
