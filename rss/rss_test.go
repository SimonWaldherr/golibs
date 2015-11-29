package rss

import (
	"testing"
)

var feeds = []string{
	"http://cre.fm/feed/m4a",
	"http://konscience.de/feed/aac/",
	"http://podcast.firtz.org/firtz/mp3",
}

func Test_Read(t *testing.T) {
	for _, url := range feeds {
		feedContent, err := Read(url)
		if err != nil {
			t.Fatalf("rss.Read Test failed: %v\nFeed: %v\n\n", err, url)
		}
		if feedContent.Title == "" {
			t.Fatalf("rss.Read Test failed: No Title detected\nFeed: %v\n\n", url)
		}
		if feedContent.Len() < 3 {
			t.Fatalf("rss.Read Test failed: I don't think that there are only %v Episodes\nFeed: %v | %v\n\n", feedContent.Len(), url, feedContent.Title)
		}

		for _, item := range feedContent.Items {
			if item.Time().Unix() < 921456000 {
				t.Fatalf("rss.Read Test failed: The episode publication date is before the release date of the first RSS draft\nFeed: %v | %v\n\n", url, feedContent.Title)
			}
		}
	}
}

func Test_Read2(t *testing.T) {
	_, err := Read("http://localhost/feed.xml")
	if err == nil {
		t.Fatal("rss.Read2 Test failed: I expect an error but got nil\n")
	}
}

func Test_Main_Time(t *testing.T) {
	feedContent, _ := Read("http://cre.fm/feed/m4a")
	if feedContent.Time().Unix() < 921456000 {
		t.Fatal("rss.Main.Time Test failed: The specified last build date is before the release date of the first RSS draft\n")
	}
}
