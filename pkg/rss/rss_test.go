package rss

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const sampleRSS = `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
<channel>
 <title>Sample RSS Feed</title>
 <description>Sample feed for testing</description>
 <link>http://www.example.com/main.html</link>
 <item>
  <title>Article 3</title>
  <description>Third article</description>
  <link>http://example.com/article3</link>
  <guid>http://example.com/article3</guid>
 </item>
 <item>
  <title>Article 2</title>
  <description>Second article</description>
  <link>http://example.com/article2</link>
  <guid>http://example.com/article2</guid>
 </item>
 <item>
  <title>Article 1</title>
  <description>First article</description>
  <link>http://example.com/article1</link>
  <guid>http://example.com/article1</guid>
 </item>
</channel>
</rss>`

func TestGetLastNews(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, sampleRSS)
	}))
	defer server.Close()

	// Test requesting 2 news items
	news, err := GetLastNews(server.URL, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(news) != 2 {
		t.Fatalf("expected 2 items, got %d", len(news))
	}

	if news[0] != "http://example.com/article3" || news[1] != "http://example.com/article2" {
		t.Errorf("unexpected news items: %v", news)
	}

	// Test requesting 0 items
	zeroNews, err := GetLastNews(server.URL, 0)
	if err != nil {
		t.Fatalf("unexpected error with 0 amount: %v", err)
	}
	if len(zeroNews) != 0 {
		t.Errorf("expected 0 items, got %d", len(zeroNews))
	}

	// Test invalid URL
	_, err = GetLastNews("http://invalid.local.nonexistent:99999", 5)
	if err == nil {
		t.Errorf("expected error with invalid URL")
	}
}
