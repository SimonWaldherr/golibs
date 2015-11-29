package rss

import (
	"testing"
)

var feeds = []string{
	"http://cre.fm/feed/m4a",
	"http://www.konscience.de/feed/aac/",
	"http://podcast.firtz.org/firtz/mp3"}

func Test_Read(t *testing.T) {
	for _, url := range feeds {
		feedContent, err := Read(url)
		if err != nil {
			t.Fatalf("rss.Read Test failed: %v\nFeed: %v\n\n", err, url)
		}
		if feedContent.Title == "" {
			t.Fatalf("rss.Read Test failed: No Title detected\nFeed:%v\n\n", url)
		}
		if feedContent.Len() < 3 {
			t.Fatalf("rss.Read Test failed: I don't think that there are only %v Episodes\nFeed:%v | %v\n\n", feedContent.Len(), url, feedContent.Title)
		}
		for _, item := range feedContent.Items {
			if item.Time().Unix() < 921456000 {
				t.Fatalf("rss.Read Test failed: The episode publication date is before the release date of the first RSS draft\nFeed:%v | %v\n\n", url, feedContent.Title)
			}
		}
	}
}
