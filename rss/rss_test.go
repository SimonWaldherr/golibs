package rss

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const testFeed = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"><channel>
<title>Example &amp; friends</title>
<link>https://example.org/</link>
<description>A local test feed</description>
<language>en</language>
<lastBuildDate>Thu, 08 Oct 2026 12:00:00 +0000</lastBuildDate>
<image><url>https://example.org/logo.png</url><title>Logo</title></image>
<item>
<title>First episode</title><link>https://example.org/1</link><guid>episode-1</guid>
<pubDate>Wed, 07 Oct 2026 10:30:00 +0000</pubDate>
<category>Science</category><category>Technology</category>
<enclosure url="https://example.org/1.mp3" type="audio/mpeg"/>
<description><![CDATA[<p>Hello world</p>]]></description>
</item>
<item><title>Second episode</title><pubDate>2026-10-06T09:00:00Z</pubDate></item>
</channel></rss>`

func assertFeed(t *testing.T, feed *Main) {
	t.Helper()
	if feed.Title != "Example & friends" || feed.Link != "https://example.org/" ||
		feed.Description != "A local test feed" || feed.Language != "en" || feed.Len() != 2 {
		t.Fatalf("unexpected feed: %+v", feed)
	}
	if feed.Image.URL != "https://example.org/logo.png" || feed.Image.Title != "Logo" {
		t.Fatalf("unexpected image: %+v", feed.Image)
	}
	item := feed.Items[0]
	if item.Title != "First episode" || item.GUID != "episode-1" ||
		item.Description != "<p>Hello world</p>" || len(item.Category) != 2 ||
		len(item.Enclosure) != 1 || item.Enclosure[0].Type != "audio/mpeg" ||
		item.Enclosure[0].URL != "https://example.org/1.mp3" {
		t.Fatalf("unexpected item: %+v", item)
	}
	if want := time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC); !feed.Time().Equal(want) {
		t.Errorf("build date = %v, want %v", feed.Time(), want)
	}
	if want := time.Date(2026, time.October, 7, 10, 30, 0, 0, time.UTC); !item.Time().Equal(want) {
		t.Errorf("publication date = %v, want %v", item.Time(), want)
	}
	if want := time.Date(2026, time.October, 6, 9, 0, 0, 0, time.UTC); !feed.Items[1].Time().Equal(want) {
		t.Errorf("RFC3339 publication date = %v, want %v", feed.Items[1].Time(), want)
	}
}

func TestParse(t *testing.T) {
	feed, err := Parse(strings.NewReader(testFeed))
	if err != nil {
		t.Fatal(err)
	}
	assertFeed(t, feed)
}

func TestRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = io.WriteString(w, testFeed)
	}))
	defer server.Close()

	feed, err := Read(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	assertFeed(t, feed)

	feed, err = ReadCustomClient(server.URL, server.Client())
	if err != nil {
		t.Fatal(err)
	}
	assertFeed(t, feed)
}

func TestParseMalformedXML(t *testing.T) {
	for _, input := range []string{"", "<rss><channel>", "<rss><channel></rss>"} {
		if feed, err := Parse(strings.NewReader(input)); err == nil || feed != nil {
			t.Errorf("Parse(%q) = %v, %v; want an error", input, feed, err)
		}
	}
}

type testClient func(string) (*http.Response, error)

func (client testClient) Get(url string) (*http.Response, error) {
	return client(url)
}

type trackedBody struct {
	io.Reader
	closed bool
}

func (body *trackedBody) Close() error {
	body.closed = true
	return nil
}

func TestReadCustomClientClosesBody(t *testing.T) {
	for _, input := range []string{testFeed, "<rss><channel>"} {
		body := &trackedBody{Reader: strings.NewReader(input)}
		client := testClient(func(url string) (*http.Response, error) {
			if url != "https://example.org/feed" {
				t.Errorf("URL = %q", url)
			}
			return &http.Response{Body: body}, nil
		})
		_, err := ReadCustomClient("https://example.org/feed", client)
		if (err == nil) != (input == testFeed) {
			t.Errorf("ReadCustomClient() error = %v", err)
		}
		if !body.closed {
			t.Error("response body was not closed")
		}
	}
}

func TestParseLeavesReaderOpen(t *testing.T) {
	body := &trackedBody{Reader: strings.NewReader(testFeed)}
	if _, err := Parse(body); err != nil {
		t.Fatal(err)
	}
	if body.closed {
		t.Error("Parse closed the caller's reader")
	}
}

func TestReadCustomClientTransportError(t *testing.T) {
	wantErr := errors.New("transport failed")
	client := testClient(func(string) (*http.Response, error) { return nil, wantErr })
	if feed, err := ReadCustomClient("https://example.org/feed", client); feed != nil || !errors.Is(err, wantErr) {
		t.Fatalf("ReadCustomClient() = %v, %v", feed, err)
	}
}

type failingReader struct{ err error }

func (reader failingReader) Read([]byte) (int, error) { return 0, reader.err }

func TestReadCustomClientBodyError(t *testing.T) {
	wantErr := errors.New("body failed")
	body := &trackedBody{Reader: failingReader{wantErr}}
	client := testClient(func(string) (*http.Response, error) {
		return &http.Response{Body: body}, nil
	})
	if feed, err := ReadCustomClient("https://example.org/feed", client); feed != nil || !errors.Is(err, wantErr) {
		t.Fatalf("ReadCustomClient() = %v, %v", feed, err)
	}
	if !body.closed {
		t.Error("response body was not closed after a read error")
	}
}
