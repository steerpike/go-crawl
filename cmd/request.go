/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crawl/models"
	"database/sql"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

func UrlExists(url string) bool {
	var exists bool
	db, err := sql.Open("sqlite3", "music.db")
	if err != nil {
		log.Println("Something went wrong opening database:", err)
	}
	query := `SELECT 1 FROM Artists WHERE Url = ?`
	err = db.QueryRow(query, url).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("querying for url existence: %v", err)
	}
	return exists
}

func CrawlURL(url string, sourceUrl string) {
	c := colly.NewCollector()
	pathSet := make(map[string]bool)
	artist := models.Artist{}
	tags := []string{}
	similar := []string{}
	videos := make(map[string]string)
	artist.SourceUrl = sourceUrl
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
		a := e.Attr("href")
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
		log.Printf("Scraped page for artist: %+v", artist)
		if strings.Contains(artist.Url, "https://www.last.fm/music") {
			artist.Save()
		}
	})
	if !UrlExists(url) {
		c.Visit(url)
	} else {
		log.Printf("URL already exists in the database: %s", url)
	}

}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch a URL for artist details and to find other related artists",
	Long: `Fetch a URL using the Colly library and find artist details,
	music, and other related artists on that page to crawl`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		CrawlURL(url, "")
	},
}

var harvestCmd = &cobra.Command{
	Use:   "harvest",
	Short: "Harvest the seed list for artist details and to find other related artists",
	Long: `Harvests a list of previously found urls to find artist details,
	music, and other related artists to add to the seed list`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", "music.db")
		if err != nil {
			log.Println("Something went wrong opening database while harvesting:", err)
		}
		rows, err := db.Query("SELECT Url, SourceUrl FROM Seeds LIMIT 1")
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			var url, sourceUrl string
			err := rows.Scan(&url, &sourceUrl)
			log.Println("url:", url)
			if err != nil {
				log.Fatal(err)
			}
			rows.Close()
			CrawlURL(url, sourceUrl)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	//fetchCmd.Flags().IntP("limit", "l", 5, "Limit of how many URLs to crawl")
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(harvestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// requestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// requestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
