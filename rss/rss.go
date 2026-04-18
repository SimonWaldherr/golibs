// Package rss provides an RSS and Atom feed parser supporting podcasts, blogs, and news.
package rss

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"simonwaldherr.de/go/golibs/as"
)

// Client is a http client to crawl the feed from (#)
type Client interface {
	Get(url string) (resp *http.Response, err error)
}

// ItemEnclosure contains an URL and its MIME-Type of an episode file
type ItemEnclosure struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Length string `xml:"length,attr"`
}

// ITunesOwner represents the iTunes podcast owner contact
type ITunesOwner struct {
	Name  string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd name"`
	Email string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd email"`
}

// ITunesCategory represents a hierarchical iTunes category
type ITunesCategory struct {
	Text        string         `xml:"text,attr"`
	SubCategory *ITunesCategory `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd category"`
}

// Item represents an episode/article with all of its content.
// It supports RSS 2.0 fields as well as iTunes podcast extensions and Dublin Core metadata.
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

	// Dublin Core metadata (http://purl.org/dc/elements/1.1/)
	DCCreator string `xml:"http://purl.org/dc/elements/1.1/ creator"`
	DCDate    string `xml:"http://purl.org/dc/elements/1.1/ date"`
	DCSubject string `xml:"http://purl.org/dc/elements/1.1/ subject"`

	// iTunes podcast extensions (http://www.itunes.com/dtds/podcast-1.0.dtd)
	ITunesTitle    string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd title"`
	ITunesDuration string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd duration"`
	ITunesExplicit string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd explicit"`
	ITunesSummary  string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd summary"`
	ITunesSubtitle string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd subtitle"`
	ITunesKeywords string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd keywords"`
	ITunesAuthor   string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd author"`
	ITunesImage    struct {
		Href string `xml:"href,attr"`
	} `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd image"`
	ITunesEpisode     string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd episode"`
	ITunesSeason      string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd season"`
	ITunesEpisodeType string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd episodeType"`
}

// Author returns the best available author name for this item,
// checking iTunes author, Dublin Core creator, and feed-level fields.
func (item *Item) Author() string {
	if item.ITunesAuthor != "" {
		return item.ITunesAuthor
	}
	if item.DCCreator != "" {
		return item.DCCreator
	}
	return ""
}

// Keywords returns the iTunes keywords as a slice of strings.
func (item *Item) Keywords() []string {
	if item.ITunesKeywords == "" {
		return nil
	}
	parts := strings.Split(item.ITunesKeywords, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// Image contains the URL and title of the podcast logo
type Image struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

// Main struct represents the root of the RSS feed.
// It supports RSS 2.0 fields, iTunes podcast extensions, and Dublin Core metadata.
type Main struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Copyright     string `xml:"copyright"`
	Description   string `xml:"description"`
	Generator     string `xml:"generator"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	PubDate       string `xml:"pubDate"`
	ManagingEditor string `xml:"managingEditor"`
	WebMaster     string `xml:"webMaster"`
	Items         []Item `xml:"item"`
	Image         Image  `xml:"image"`

	// Dublin Core metadata (http://purl.org/dc/elements/1.1/)
	DCCreator string `xml:"http://purl.org/dc/elements/1.1/ creator"`

	// iTunes podcast extensions (http://www.itunes.com/dtds/podcast-1.0.dtd)
	ITunesAuthor   string         `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd author"`
	ITunesOwner    ITunesOwner    `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd owner"`
	ITunesSummary  string         `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd summary"`
	ITunesSubtitle string         `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd subtitle"`
	ITunesExplicit string         `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd explicit"`
	ITunesKeywords string         `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd keywords"`
	ITunesCategory ITunesCategory `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd category"`
	ITunesImage    struct {
		Href string `xml:"href,attr"`
	} `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd image"`
	ITunesType       string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd type"`
	ITunesNewFeedURL string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd new-feed-url"`
}

// IsPodcast returns true if the feed contains iTunes podcast metadata.
func (rss *Main) IsPodcast() bool {
	if rss == nil {
		return false
	}
	return rss.ITunesAuthor != "" || rss.ITunesOwner.Email != "" || rss.ITunesCategory.Text != ""
}

// Author returns the best available author/editor for the feed.
func (rss *Main) Author() string {
	if rss == nil {
		return ""
	}
	if rss.ITunesAuthor != "" {
		return rss.ITunesAuthor
	}
	if rss.ManagingEditor != "" {
		return rss.ManagingEditor
	}
	if rss.DCCreator != "" {
		return rss.DCCreator
	}
	return ""
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

// AtomLink represents an Atom link element.
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

// AtomPerson represents an Atom person (author or contributor).
type AtomPerson struct {
	Name  string `xml:"name"`
	Email string `xml:"email"`
	URI   string `xml:"uri"`
}

// AtomEntry represents a single entry in an Atom feed (blog post, news article, etc.).
type AtomEntry struct {
	ID        string       `xml:"id"`
	Title     string       `xml:"title"`
	Updated   string       `xml:"updated"`
	Published string       `xml:"published"`
	Authors   []AtomPerson `xml:"author"`
	Links     []AtomLink   `xml:"link"`
	Summary   string       `xml:"summary"`
	Content   string       `xml:"content"`
	Category  []struct {
		Term  string `xml:"term,attr"`
		Label string `xml:"label,attr"`
	} `xml:"category"`
	Rights string `xml:"rights"`
}

// Link returns the first alternate HTML link, or the first link if no alternate found.
func (e *AtomEntry) Link() string {
	for _, l := range e.Links {
		if l.Rel == "alternate" || l.Rel == "" {
			return l.Href
		}
	}
	if len(e.Links) > 0 {
		return e.Links[0].Href
	}
	return ""
}

// Time returns the publication time of the entry.
func (e *AtomEntry) Time() time.Time {
	if e.Published != "" {
		return as.Time(e.Published)
	}
	return as.Time(e.Updated)
}

// Author returns the primary author name of the entry.
func (e *AtomEntry) Author() string {
	if len(e.Authors) > 0 {
		return e.Authors[0].Name
	}
	return ""
}

// AtomFeed represents an Atom 1.0 feed (RFC 4287).
// Atom is widely used for blogs and news sites.
type AtomFeed struct {
	ID        string       `xml:"id"`
	Title     string       `xml:"title"`
	Subtitle  string       `xml:"subtitle"`
	Updated   string       `xml:"updated"`
	Authors   []AtomPerson `xml:"author"`
	Links     []AtomLink   `xml:"link"`
	Rights    string       `xml:"rights"`
	Generator string       `xml:"generator"`
	Entries   []AtomEntry  `xml:"entry"`
}

// Len returns the number of entries in the Atom feed.
func (f *AtomFeed) Len() int {
	if f == nil {
		return 0
	}
	return len(f.Entries)
}

// Time returns the last updated time of the feed.
func (f *AtomFeed) Time() time.Time {
	if f == nil {
		return time.Time{}
	}
	return as.Time(f.Updated)
}

// Link returns the primary HTML link for the feed.
func (f *AtomFeed) Link() string {
	if f == nil {
		return ""
	}
	for _, l := range f.Links {
		if l.Rel == "alternate" || l.Rel == "" {
			return l.Href
		}
	}
	if len(f.Links) > 0 {
		return f.Links[0].Href
	}
	return ""
}

// Author returns the primary author name of the feed.
func (f *AtomFeed) Author() string {
	if f == nil || len(f.Authors) == 0 {
		return ""
	}
	return f.Authors[0].Name
}

// ReadAtom fetches an Atom feed from the given URL.
func ReadAtom(url string) (*AtomFeed, error) {
	return ReadAtomCustomClient(url, http.DefaultClient)
}

// ReadAtomCustomClient fetches an Atom feed using a custom HTTP client.
func ReadAtomCustomClient(url string, client Client) (*AtomFeed, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var feed AtomFeed
	if err := xml.NewDecoder(r.Body).Decode(&feed); err != nil {
		return nil, err
	}
	return &feed, nil
}

// FeedType describes the detected type of a feed.
type FeedType int

const (
	FeedTypeUnknown FeedType = iota
	FeedTypeRSS
	FeedTypeAtom
)

// Feed is a unified representation of either an RSS or Atom feed.
type Feed struct {
	Type FeedType

	// RSS is set when Type == FeedTypeRSS
	RSS *Main

	// Atom is set when Type == FeedTypeAtom
	Atom *AtomFeed
}

// Title returns the feed title regardless of type.
func (f *Feed) Title() string {
	switch f.Type {
	case FeedTypeRSS:
		return f.RSS.Title
	case FeedTypeAtom:
		return f.Atom.Title
	}
	return ""
}

// Len returns the number of items/entries regardless of type.
func (f *Feed) Len() int {
	switch f.Type {
	case FeedTypeRSS:
		return f.RSS.Len()
	case FeedTypeAtom:
		return f.Atom.Len()
	}
	return 0
}

// ReadAuto fetches a URL and automatically detects whether it is an RSS or Atom feed.
func ReadAuto(url string) (*Feed, error) {
	return ReadAutoCustomClient(url, http.DefaultClient)
}

// ReadAutoCustomClient fetches a URL and detects the feed type using a custom HTTP client.
// It supports RSS 2.0, RSS 1.0 (RDF), and Atom 1.0.
func ReadAutoCustomClient(url string, client Client) (*Feed, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// Scan the token stream to find the first start element, then detect feed type.
	dec := xml.NewDecoder(r.Body)
	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		se, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}

		local := strings.ToLower(se.Name.Local)
		switch local {
		case "rss", "rdf":
			// RSS 2.0 root is <rss>, RSS 1.0 (RDF) root is <rdf:RDF> (Local="RDF")
			var wrapper struct {
				Main Main `xml:"channel"`
			}
			if err := dec.DecodeElement(&wrapper, &se); err != nil {
				return nil, err
			}
			return &Feed{Type: FeedTypeRSS, RSS: &wrapper.Main}, nil
		case "feed":
			// Atom 1.0 root is <feed>
			var feed AtomFeed
			if err := dec.DecodeElement(&feed, &se); err != nil {
				return nil, err
			}
			return &Feed{Type: FeedTypeAtom, Atom: &feed}, nil
		default:
			return nil, fmt.Errorf("unsupported feed root element: %s", se.Name.Local)
		}
	}
}

// Len returns the episode count
func (rss *Main) Len() int {
	if rss == nil {
		return 0
	}
	return len(rss.Items)
}

// Time returns the specified last build date as time.Time
func (rss *Main) Time() time.Time {
	if rss == nil {
		return time.Time{}
	}
	if rss.LastBuildDate != "" {
		return as.Time(rss.LastBuildDate)
	}
	return as.Time(rss.PubDate)
}

// Time returns the specified publication date as time.Time
func (episode *Item) Time() time.Time {
	if episode == nil {
		return time.Time{}
	}
	return as.Time(episode.PubDate)
}
