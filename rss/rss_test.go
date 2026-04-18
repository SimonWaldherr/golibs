package rss

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// sampleRSS is a minimal RSS 2.0 podcast feed for local tests.
const sampleRSS = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
  xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd"
  xmlns:dc="http://purl.org/dc/elements/1.1/">
  <channel>
    <title>Test Podcast</title>
    <link>https://example.com/podcast</link>
    <description>A test podcast feed</description>
    <language>en</language>
    <lastBuildDate>Thu, 01 Jan 2015 00:00:00 +0000</lastBuildDate>
    <itunes:author>Test Author</itunes:author>
    <itunes:category text="Technology"/>
    <item>
      <title>Episode 1</title>
      <link>https://example.com/ep1</link>
      <pubDate>Mon, 01 Jun 2020 00:00:00 +0000</pubDate>
      <guid>https://example.com/ep1</guid>
      <dc:creator>Test Author</dc:creator>
      <itunes:duration>01:23:45</itunes:duration>
      <itunes:explicit>no</itunes:explicit>
      <itunes:keywords>test, podcast, episode</itunes:keywords>
      <enclosure url="https://example.com/ep1.mp3" type="audio/mpeg" length="12345678"/>
      <description>First episode</description>
    </item>
    <item>
      <title>Episode 2</title>
      <link>https://example.com/ep2</link>
      <pubDate>Tue, 02 Jun 2020 00:00:00 +0000</pubDate>
      <guid>https://example.com/ep2</guid>
      <itunes:duration>00:45:00</itunes:duration>
      <enclosure url="https://example.com/ep2.mp3" type="audio/mpeg" length="9876543"/>
      <description>Second episode</description>
    </item>
    <item>
      <title>Episode 3</title>
      <link>https://example.com/ep3</link>
      <pubDate>Wed, 03 Jun 2020 00:00:00 +0000</pubDate>
      <guid>https://example.com/ep3</guid>
      <enclosure url="https://example.com/ep3.mp3" type="audio/mpeg" length="7654321"/>
      <description>Third episode</description>
    </item>
  </channel>
</rss>`

// sampleAtom is a minimal Atom 1.0 feed for local tests.
const sampleAtom = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <id>https://example.com/blog</id>
  <title>Test Blog</title>
  <subtitle>A test Atom feed</subtitle>
  <updated>2021-06-01T00:00:00Z</updated>
  <author><name>Blog Author</name><email>author@example.com</email></author>
  <link href="https://example.com/blog" rel="alternate"/>
  <link href="https://example.com/atom.xml" rel="self"/>
  <entry>
    <id>https://example.com/blog/post1</id>
    <title>Post 1</title>
    <published>2021-05-01T10:00:00Z</published>
    <updated>2021-05-01T10:00:00Z</updated>
    <author><name>Blog Author</name></author>
    <link href="https://example.com/blog/post1" rel="alternate"/>
    <summary>First post summary</summary>
    <content>First post content</content>
  </entry>
  <entry>
    <id>https://example.com/blog/post2</id>
    <title>Post 2</title>
    <published>2021-06-01T10:00:00Z</published>
    <updated>2021-06-01T10:00:00Z</updated>
    <link href="https://example.com/blog/post2" rel="alternate"/>
    <summary>Second post summary</summary>
  </entry>
</feed>`

func newMockServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Write([]byte(body))
	}))
}

func Test_Read_RSS(t *testing.T) {
	ts := newMockServer(sampleRSS)
	defer ts.Close()

	feed, err := Read(ts.URL)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if feed.Title != "Test Podcast" {
		t.Errorf("expected title %q, got %q", "Test Podcast", feed.Title)
	}
	if feed.Len() != 3 {
		t.Errorf("expected 3 items, got %d", feed.Len())
	}
	if !feed.IsPodcast() {
		t.Errorf("expected IsPodcast() == true")
	}
	if feed.Author() != "Test Author" {
		t.Errorf("expected author %q, got %q", "Test Author", feed.Author())
	}
}

func Test_Read_RSSTime(t *testing.T) {
	ts := newMockServer(sampleRSS)
	defer ts.Close()

	feed, err := Read(ts.URL)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	ft := feed.Time()
	if ft.IsZero() {
		t.Error("expected non-zero time from feed")
	}
}

func Test_ItemTime(t *testing.T) {
	ts := newMockServer(sampleRSS)
	defer ts.Close()

	feed, err := Read(ts.URL)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	for _, item := range feed.Items {
		it := item.Time()
		if it.IsZero() {
			t.Errorf("expected non-zero time for item %q", item.Title)
		}
	}
}

func Test_ItemAuthorAndKeywords(t *testing.T) {
	ts := newMockServer(sampleRSS)
	defer ts.Close()

	feed, err := Read(ts.URL)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	ep1 := feed.Items[0]
	if ep1.Author() != "Test Author" {
		t.Errorf("expected item author %q, got %q", "Test Author", ep1.Author())
	}
	kw := ep1.Keywords()
	if len(kw) != 3 {
		t.Errorf("expected 3 keywords, got %d: %v", len(kw), kw)
	}
}

func Test_Read_Atom(t *testing.T) {
	ts := newMockServer(sampleAtom)
	defer ts.Close()

	feed, err := ReadAtom(ts.URL)
	if err != nil {
		t.Fatalf("ReadAtom failed: %v", err)
	}
	if feed.Title != "Test Blog" {
		t.Errorf("expected title %q, got %q", "Test Blog", feed.Title)
	}
	if feed.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", feed.Len())
	}
	if feed.Author() != "Blog Author" {
		t.Errorf("expected author %q, got %q", "Blog Author", feed.Author())
	}
	if feed.Link() != "https://example.com/blog" {
		t.Errorf("expected link %q, got %q", "https://example.com/blog", feed.Link())
	}
}

func Test_AtomEntryFields(t *testing.T) {
	ts := newMockServer(sampleAtom)
	defer ts.Close()

	feed, err := ReadAtom(ts.URL)
	if err != nil {
		t.Fatalf("ReadAtom failed: %v", err)
	}
	e := feed.Entries[0]
	if e.Title != "Post 1" {
		t.Errorf("expected title %q, got %q", "Post 1", e.Title)
	}
	if e.Link() != "https://example.com/blog/post1" {
		t.Errorf("expected link %q, got %q", "https://example.com/blog/post1", e.Link())
	}
	if e.Time().IsZero() {
		t.Error("expected non-zero entry time")
	}
	if e.Author() != "Blog Author" {
		t.Errorf("expected author %q, got %q", "Blog Author", e.Author())
	}
}

func Test_ReadAuto_RSS(t *testing.T) {
	ts := newMockServer(sampleRSS)
	defer ts.Close()

	feed, err := ReadAuto(ts.URL)
	if err != nil {
		t.Fatalf("ReadAuto failed: %v", err)
	}
	if feed.Type != FeedTypeRSS {
		t.Errorf("expected FeedTypeRSS, got %v", feed.Type)
	}
	if feed.Title() != "Test Podcast" {
		t.Errorf("expected title %q, got %q", "Test Podcast", feed.Title())
	}
	if feed.Len() != 3 {
		t.Errorf("expected 3 items, got %d", feed.Len())
	}
}

func Test_ReadAuto_Atom(t *testing.T) {
	ts := newMockServer(sampleAtom)
	defer ts.Close()

	feed, err := ReadAuto(ts.URL)
	if err != nil {
		t.Fatalf("ReadAuto failed: %v", err)
	}
	if feed.Type != FeedTypeAtom {
		t.Errorf("expected FeedTypeAtom, got %v", feed.Type)
	}
	if feed.Title() != "Test Blog" {
		t.Errorf("expected title %q, got %q", "Test Blog", feed.Title())
	}
	if feed.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", feed.Len())
	}
}

func Test_Read_Error(t *testing.T) {
	_, err := Read("http://127.0.0.1:1/nonexistent")
	if err == nil {
		t.Error("expected error for unreachable URL, got nil")
	}
}

func Test_Read_InvalidXML(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not valid xml"))
	}))
	defer ts.Close()

	_, err := Read(ts.URL)
	if err == nil {
		t.Error("expected error for invalid XML, got nil")
	}
}

func Test_CustomClient(t *testing.T) {
	ts := newMockServer(sampleRSS)
	defer ts.Close()

	feed, err := ReadCustomClient(ts.URL, http.DefaultClient)
	if err != nil {
		t.Fatalf("ReadCustomClient failed: %v", err)
	}
	if feed.Len() < 1 {
		t.Error("expected at least one item")
	}
}

// Test_Read_Network exercises real external RSS feeds.
// Skipped when testing.Short() is true (e.g. in CI without network access).
func Test_Read_Network(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	// A small curated list of stable, long-running feeds.
	feeds := []string{
		"https://feeds.feedburner.com/TheTwimlAiPodcast",
		"https://changelog.com/gotime/feed",
	}

	for _, url := range feeds {
		url := url
		t.Run(url, func(t *testing.T) {
			t.Parallel()
			feed, err := Read(url)
			if err != nil {
				t.Errorf("Read(%q) failed: %v", url, err)
				return
			}
			if feed.Title == "" {
				t.Errorf("Read(%q): empty title", url)
			}
			if feed.Len() < 1 {
				t.Errorf("Read(%q): no items", url)
			}
		})
	}
}

// Test_NilReceiver ensures nil-safe methods don't panic.
func Test_NilReceiver(t *testing.T) {
	var m *Main
	if m.Time().IsZero() == false {
		t.Error("nil Main.Time() should return zero time")
	}
	if m.Len() != 0 {
		t.Error("nil Main.Len() should return 0")
	}
	if m.Author() != "" {
		t.Error("nil Main.Author() should return empty string")
	}
	if m.IsPodcast() != false {
		t.Error("nil Main.IsPodcast() should return false")
	}

	var e *Item
	if e.Time().IsZero() == false {
		t.Error("nil Item.Time() should return zero time")
	}

	var f *AtomFeed
	if f.Time().IsZero() == false {
		t.Error("nil AtomFeed.Time() should return zero time")
	}
	if f.Len() != 0 {
		t.Error("nil AtomFeed.Len() should return 0")
	}
}

// Test_ReadAtomCustomClient verifies that custom clients work for Atom feeds.
func Test_ReadAtomCustomClient(t *testing.T) {
	ts := newMockServer(sampleAtom)
	defer ts.Close()

	feed, err := ReadAtomCustomClient(ts.URL, http.DefaultClient)
	if err != nil {
		t.Fatalf("ReadAtomCustomClient failed: %v", err)
	}
	if feed.Len() < 1 {
		t.Error("expected at least one entry")
	}
}

// Test_EnclosureFields verifies that enclosure metadata is parsed correctly.
func Test_EnclosureFields(t *testing.T) {
	ts := newMockServer(sampleRSS)
	defer ts.Close()

	feed, err := Read(ts.URL)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if len(feed.Items[0].Enclosure) == 0 {
		t.Fatal("expected at least one enclosure in item 0")
	}
	enc := feed.Items[0].Enclosure[0]
	if !strings.HasPrefix(enc.URL, "https://") {
		t.Errorf("expected enclosure URL to start with https://, got %q", enc.URL)
	}
	if enc.Type != "audio/mpeg" {
		t.Errorf("expected enclosure type %q, got %q", "audio/mpeg", enc.Type)
	}
	if enc.Length != 12345678 {
		t.Errorf("expected enclosure length %d, got %d", int64(12345678), enc.Length)
	}
}

// Test_ReadAuto_RootElement_RDF tests that RDF-style RSS feeds (RSS 1.0) are recognized.
// Note: RSS 1.0 / RDF feeds use a flat structure where items are siblings of <channel>,
// not nested inside it, so only channel-level metadata can be parsed, not items.
func Test_ReadAuto_RootElement_RDF(t *testing.T) {
	rdfFeed := `<?xml version="1.0" encoding="UTF-8"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
  xmlns="http://purl.org/rss/1.0/">
  <channel rdf:about="https://example.com/feed">
    <title>RDF Feed</title>
    <link>https://example.com</link>
    <description>An RDF/RSS 1.0 feed</description>
  </channel>
  <item rdf:about="https://example.com/1">
    <title>Item 1</title>
    <link>https://example.com/1</link>
  </item>
</rdf:RDF>`

	ts := newMockServer(rdfFeed)
	defer ts.Close()

	// RDF feeds start with <rdf:RDF> (Local="RDF"), lowercased to "rdf" → treated as FeedTypeRSS.
	feed, err := ReadAuto(ts.URL)
	if err != nil {
		t.Fatalf("ReadAuto RDF failed: %v", err)
	}
	if feed.Type != FeedTypeRSS {
		t.Errorf("expected FeedTypeRSS for RDF feed, got %v", feed.Type)
	}
	// Channel-level metadata is parsed
	if feed.RSS.Title != "RDF Feed" {
		t.Errorf("expected title %q, got %q", "RDF Feed", feed.RSS.Title)
	}
}

// Verify that xml.Unmarshal recognises the iTunes namespace correctly.
func Test_ITunesFields(t *testing.T) {
	raw := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
  <channel>
    <title>Podcast</title>
    <itunes:author>Jane Doe</itunes:author>
    <itunes:owner>
      <itunes:name>Jane Doe</itunes:name>
      <itunes:email>jane@example.com</itunes:email>
    </itunes:owner>
    <itunes:category text="Science"/>
    <item>
      <title>Ep 1</title>
      <pubDate>Mon, 01 Jun 2020 00:00:00 +0000</pubDate>
      <itunes:duration>00:30:00</itunes:duration>
      <itunes:episode>1</itunes:episode>
      <itunes:season>1</itunes:season>
    </item>
  </channel>
</rss>`

	var wrapper struct {
		Main Main `xml:"channel"`
	}
	if err := xml.Unmarshal([]byte(raw), &wrapper); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	m := wrapper.Main
	if m.ITunesAuthor != "Jane Doe" {
		t.Errorf("expected ITunesAuthor %q, got %q", "Jane Doe", m.ITunesAuthor)
	}
	if m.ITunesOwner.Email != "jane@example.com" {
		t.Errorf("expected owner email %q, got %q", "jane@example.com", m.ITunesOwner.Email)
	}
	if m.ITunesCategory.Text != "Science" {
		t.Errorf("expected category %q, got %q", "Science", m.ITunesCategory.Text)
	}
	ep := m.Items[0]
	if ep.ITunesDuration != "00:30:00" {
		t.Errorf("expected duration %q, got %q", "00:30:00", ep.ITunesDuration)
	}
	if ep.ITunesEpisode != "1" {
		t.Errorf("expected episode %q, got %q", "1", ep.ITunesEpisode)
	}
}
