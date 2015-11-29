//Package rss provides a RSS parser
package rss

import (
	"encoding/xml"
	"net/http"
	"simonwaldherr.de/go/golibs/as"
	"time"
)

// Client is a http client to crawl the feed from (#)
type Client interface {
	Get(url string) (resp *http.Response, err error)
}

// ItemEnclosure contains an URL and its MIME-Type of an episode file
type ItemEnclosure struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

// Item represents an episode with all of its content
type Item struct {
	Title       string          `xml:"title"`
	Link        string          `xml:"link"`
	Comments    string          `xml:"comments"`
	PubDate     string          `xml:"pubDate"`
	GUID        string          `xml:"guid"`
	Category    []string        `xml:"category"`
	Enclosure   []ItemEnclosure `xml:"enclosure"`
	Description string          `xml:"description"`
	Content     string          `xml:"content"`
}

// Image contains the URL and title of the podcast logo
type Image struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
}

// Main struct represents the root of the RSS feed
type Main struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Copyright     string `xml:"copyright"`
	Description   string `xml:"description"`
	Generator     string `xml:"generator"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Items         []Item `xml:"item"`
	Image         Image  `xml:"image"`
}

// Read reads from an URL with the Golang default http client and returns rss.Main
func Read(url string) (*Main, error) {
	return ReadCustomClient(url, http.DefaultClient)
}

// ReadCustomClient allows to read with an custom http client
func ReadCustomClient(url string, client Client) (*Main, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r.Body.Close()
	}()
	xmlDecoder := xml.NewDecoder(r.Body)

	var rss struct {
		Main Main `xml:"channel"`
	}
	if err = xmlDecoder.Decode(&rss); err != nil {
		return nil, err
	}
	return &rss.Main, nil
}

// Len returns the episode count
func (rss *Main) Len() int {
	return len(rss.Items)
}

// Time returns the specified last build date as time.Time
func (rss *Main) Time() time.Time {
	return as.Time(rss.LastBuildDate)
}

// Time returns the specified publication date as time.Time
func (episode *Item) Time() time.Time {
	return as.Time(episode.PubDate)
}
