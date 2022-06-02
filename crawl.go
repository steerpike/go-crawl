package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/gocolly/colly"
)

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
}

func (ac ArtistCrawler) GetArtistUrl(e *colly.HTMLElement) string {
	return e.ChildAttr(ac.UrlTag, ac.UrlAttribute)
}

func main() {
	target := os.Args[1]
	if IsValidUrl(target) {
		startCrawl(target)
	} else {
		fmt.Println("Please provide a lastfm artist url")
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

func startCrawl(url string) {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	c.OnHTML(`html`, func(e *colly.HTMLElement) {
		ac := ArtistCrawler{JsonTag: "#tlmdata", JsonAttribute: "data-tealium-data", UrlTag: "link[rel=canonical]", UrlAttribute: "href"}
		log.Println(e.ChildAttr(ac.JsonTag, ac.JsonAttribute))
		fmt.Println(e.ChildAttr("link[rel=canonical]", "href"))
		e.ForEach("li[class=tag]", func(_ int, kf *colly.HTMLElement) {
			fmt.Println(kf.ChildAttr("a", "href"))
			fmt.Println(kf.ChildText("a"))
		})
		e.ForEach("h3[class=artist-similar-artists-sidebar-item-name]", func(_ int, kf *colly.HTMLElement) {
			fmt.Println(kf.ChildAttr("a", "href"))
			fmt.Println(kf.ChildText("a"))
		})
		e.ForEach("td[class=chartlist-play]", func(_ int, x *colly.HTMLElement) {
			fmt.Println(x.ChildAttr("a", "href"))
			fmt.Println(x.ChildAttr("a", "data-track-name"))
		})
	})

	c.Visit(url)
}
