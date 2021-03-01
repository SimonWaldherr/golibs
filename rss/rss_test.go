package rss

import (
	"sync"
	"testing"
	"time"
)

var feeds = []string{
	"http://cre.fm/feed/m4a",
	//"http://www.konscience.de/feed/",
	//"http://minkorrekt.de/feed/m4a/",
	"http://modellansatz.de/rss/?psb",
	"http://forschergeist.de/feed/m4a/",
	//"http://www.radiomono.net/feed/m4a/",
	"http://www.zeitsprung.fm/feed/mp4/",
	//"http://podcast.firtz.org/firtz/mp3",
	"http://www.psycho-talk.de/feed/mp3/",
	"http://www.exponiert.berlin/feed/m4a/",
	"http://resonator-podcast.de/feed/m4a/",
	"http://feeds.metaebene.me/raumzeit/m4a",
	"http://www.openscienceradio.de/feed/mp4/",
	//"http://www.hoaxilla.de/podcast/hoaxilla.xml",
	"http://omegataupodcast.net/category/podcast/feed/",
	"http://www.uibk.ac.at/downloads/c115/zeit/zeit_mp4.xml",
}

func initTimeSeq() func() int {
	t := time.Now().Unix()
	return func() int {
		return int(time.Now().Unix() - t)
	}
}

func Test_Read(t *testing.T) {
	var feedsContent = make(map[string]*Main)
	var wg sync.WaitGroup

	timeSince := initTimeSeq()

	for _, url := range feeds {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			feedContent, err := Read(url)
			t.Logf("%v is ready (%v) ... \n", url, timeSince())
			if err == nil {
				feedsContent[url] = feedContent
			} else {
				t.Fatalf("rss.Read Test failed: %v\nFeed: %v\n\n", err, url)
			}
		}(url)
	}

	wg.Wait()

	for url, feedContent := range feedsContent {
		t.Logf("URL: %v\tTitle: %v\n", url, feedContent.Title)
		if feedContent.Title == "" {
			t.Fatalf("rss.Read Test failed: No Title detected\nFeed: %v\n\n", url)
		}
		if feedContent.Len() < 3 {
			t.Fatalf("rss.Read Test failed: I don't think that there are only %v Episodes\nFeed: %v | %v\n\n", feedContent.Len(), url, feedContent.Title)
		}

		for _, item := range feedContent.Items {
			t.Logf("Title: %v\tTime: %v\n", item.Title, item.Time())
			if item.Time().Unix() < 921456000 && url != "http://modellansatz.de/rss/?psb" {
				t.Fatalf("rss.Read Test failed: The episode publication date is before the release date of the first RSS draft\nFeed: %v | %v \nTime: %v\nTitle: %v\n\n", url, feedContent.Title, item.Time(), item.Title)
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
