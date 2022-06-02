package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/artist", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>Gordi</title>
<link rel="canonical" href="https://www.last.fm/music/Gordi" data-replaceable-head-tag />
</head>
<body>
<div
    id="initial-tealium-data"
    data-require="tracking/tealium-utag-set"
    data-tealium-data="{&#34;siteSection&#34;: &#34;music&#34;, &#34;pageType&#34;: &#34;artist_door&#34;, &#34;pageName&#34;: &#34;music/artist/overview&#34;, &#34;nativeEventTracking&#34;: true, &#34;userState&#34;: &#34;not authenticated&#34;, &#34;userType&#34;: &#34;anon&#34;, &#34;musicArtistName&#34;: &#34;Gordi&#34;, &#34;artist&#34;: &#34;gordi&#34;, &#34;ar&#34;: &#34;gordi,wesleygonzalez,sad13,loyallobos,whitewizzard&#34;, &#34;tag&#34;: &#34;indiepop,indie,heavymetal,femalevocalists,australia&#34;}"
    data-tealium-environment="prod"></div>

<h1>Hello World</h1>
<p class="description">This is a test page</p>
<p class="description">This is a test paragraph</p>

<ul class="
   tags-list
   tags-list--global
   ">
   <li
      class="tag"
      ><a
      href="/tag/indie+pop"
      >indie pop</a></li>
   <li
      class="tag"
      ><a
      href="/tag/indie"
      >indie</a></li>
   <li
      class="tag"
      ><a
      href="/tag/heavy+metal"
      >heavy metal</a></li>
   <li
      class="tag"
      ><a
      href="/tag/female+vocalists"
      >female vocalists</a></li>
   <li
      class="tag"
      ><a
      href="/tag/australia"
      >australia</a></li>
</ul>
<h3 class="artist-similar-artists-sidebar-item-name" itemprop="name">
	<a href="/music/Wesley+Gonzalez" itemprop="url" class="link-block-target">Wesley Gonzalez</a>
</h3>
<h3 class="artist-similar-artists-sidebar-item-name" itemprop="name">
	<a href="/music/sad13" itemprop="url" class="link-block-target">sad13</a>
</h3>
<h3 class="artist-similar-artists-sidebar-item-name" itemprop="name">
	<a href="/music/Loyal+Lobos" itemprop="url" class="link-block-target">Loyal Lobos</a>
</h3>
</body>
</html>
		`))
	})
	return httptest.NewServer(mux)
}

func TestValidateFlags(t *testing.T) {
	t.Run("request is a valid artist url", func(t *testing.T) {
		url := "https://www.last.fm/music/Gordi/"
		assert.Equal(t, true, IsValidUrl(url))
	})
	t.Run("request is a lastfm url but invalid", func(t *testing.T) {
		url := "https://www.last.fm/music"
		assert.Equal(t, false, IsValidUrl(url))
	})
	t.Run("request is a lastfm url but too long", func(t *testing.T) {
		url := "https://www.last.fm/music/Kate+Bush/_/Running+Up+That+Hill+(A+Deal+with+God)"
		assert.Equal(t, false, IsValidUrl(url))
	})
	t.Run("request is an invalid url", func(t *testing.T) {
		url := "https://www.google.com"
		assert.Equal(t, false, IsValidUrl(url))
	})
}

func TestCollectorOnArtist(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()
	c := colly.NewCollector()
	ac := ArtistCrawler{
		JsonTag: "#tlmdata", JsonAttribute: "data-tealium-data",
		UrlTag: "link[rel=canonical]", UrlAttribute: "href",
	}
	titleCallbackCalled := false
	c.OnHTML("title", func(e *colly.HTMLElement) {
		titleCallbackCalled = true
		if e.Text != "Gordi" {
			t.Error("Title element text does not match, got", e.Text)
		}
	})
	c.OnHTML("html", func(e *colly.HTMLElement) {
		artistUrl := ac.GetArtistUrl(e)
		if artistUrl != "https://www.last.fm/music/Gordi" {
			t.Error("Found incorrect artist url, got", artistUrl)
		}
	})
	c.Visit(ts.URL + "/artist")

	if !titleCallbackCalled {
		t.Error("Failed to call OnHTML callback for <title> tag")
	}
}