package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
)

const fileName = "music.db"

type ArtistCrawler struct {
	JsonTag                    string
	JsonAttribute              string
	UrlTag                     string
	UrlAttribute               string
	TagsAggregateTag           string
	TagLinkTag                 string
	TagNameAttribute           string
	SimilarArtistsAggregateTag string
	SimilarArtistsTag          string
	SimilarArtistsAttribute    string
	VideosAggregateTag         string
	VideoTag                   string
	VideoAttribute             string
	VideoNameAttribute         string
}

func (ac ArtistCrawler) GetArtistUrl(e *colly.HTMLElement) string {
	return e.ChildAttr(ac.UrlTag, ac.UrlAttribute)
}

func (ac ArtistCrawler) GetArtistJsonString(e *colly.HTMLElement) string {
	return e.ChildAttr(ac.JsonTag, ac.JsonAttribute)
}

func (ac ArtistCrawler) GetArtistNameFromJson(e *colly.HTMLElement) string {
	text := ac.GetArtistJsonString(e)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		log.Println(err)
		return ""
	}
	name := fmt.Sprintf("%v", result["musicArtistName"])
	return name
}

func (ac ArtistCrawler) GetArtistTags(e *colly.HTMLElement) map[string]string {
	tags := map[string]string{}
	e.ForEach(ac.TagsAggregateTag, func(_ int, kf *colly.HTMLElement) {
		tagPath := kf.ChildAttr(ac.TagLinkTag, ac.TagNameAttribute)
		tagName := kf.ChildText(ac.TagLinkTag)
		tags[tagPath] = tagName
	})
	return tags
}

func (ac ArtistCrawler) GetSimilarArtists(e *colly.HTMLElement) map[string]string {
	similarArtists := map[string]string{}
	e.ForEach(ac.SimilarArtistsAggregateTag, func(_ int, kf *colly.HTMLElement) {
		artistPath := kf.ChildAttr(ac.SimilarArtistsTag, ac.SimilarArtistsAttribute)
		artistName := kf.ChildText(ac.SimilarArtistsTag)
		similarArtists[artistPath] = artistName
	})
	return similarArtists
}

func (ac ArtistCrawler) GetVideoLinks(e *colly.HTMLElement) map[string]string {
	videos := map[string]string{}
	e.ForEach(ac.VideosAggregateTag, func(_ int, kf *colly.HTMLElement) {
		videoPath := kf.ChildAttr(ac.VideoTag, ac.VideoAttribute)
		videoName := kf.ChildAttr(ac.VideoTag, ac.VideoNameAttribute)
		log.Println(videoPath)
		log.Println(videoName)
		videos[videoPath] = videoName
	})
	return videos
}

func main() {
	if len(os.Args) > 1 && IsValidUrl(os.Args[1]) {
		startCrawl(os.Args[1])
	} else {
		log.Println("Please provide a lastfm artist url.")
	}

}
func IsValidUrl(input string) bool {
	u, err := url.Parse(input)
	if err != nil {
		log.Fatal(err)
		return false
	}
	if u.Host == "" || !IsLastFMHost(u.Host) {
		return false
	}
	if !IsLastFMArtistPath(u.Path) {
		return false
	}
	return true
}

func IsLastFMHost(input string) bool {
	allowed := []string{"last.fm", "www.last.fm"}
	return slices.Contains(allowed, input)
}

func IsLastFMArtistPath(input string) bool {

	splitFn := func(c rune) bool {
		return c == '/'
	}
	paths := strings.FieldsFunc(input, splitFn)
	if len(paths) != 2 {
		return false
	}
	return slices.Contains(paths, "music")
}

func startCrawl(website string) {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println(r.StatusCode)
	})

	c.OnHTML(`html`, func(e *colly.HTMLElement) {
		ac := ArtistCrawler{
			JsonTag: "#tlmdata", JsonAttribute: "data-tealium-data",
			UrlTag: "link[rel=canonical]", UrlAttribute: "href",
			TagsAggregateTag: "li[class=tag]", TagLinkTag: "a",
			TagNameAttribute:           "href",
			SimilarArtistsAggregateTag: "h3[class=artist-similar-artists-sidebar-item-name]",
			SimilarArtistsTag:          "a", SimilarArtistsAttribute: "href",
			VideosAggregateTag: "td[class=chartlist-play]", VideoTag: "a",
			VideoAttribute: "href", VideoNameAttribute: "data-track-name",
		}
		artist := Artist{}
		artist.Name = ac.GetArtistNameFromJson(e)
		artist.Url = ac.GetArtistUrl(e)
		u, err := url.Parse(artist.Url)
		if err != nil {
			panic(err)
		}
		artist.Path = u.Path
		ac.GetArtistTags(e)
		ac.GetSimilarArtists(e)
		ac.GetVideoLinks(e)
		storeArtist(artist)
	})
	c.Visit(website)
}

func storeArtist(a Artist) *Artist {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}
	repository := NewSQLiteRespository(db)
	er := repository.InitTables()
	if er != nil {
		log.Fatal("Creating table", er)
	}
	artist, err := repository.Insert(a)
	return artist
}
