package pastebin

import (
	"strconv"
	"testing"
	"time"
)

func TestJsonPasteToPaste(t *testing.T) {
	testCases := []struct {
		desc      string
		jsonPaste *jsonPaste

		assert func(*testing.T, *Paste)
	}{
		{
			desc:      "empty json paste to paste works",
			jsonPaste: &jsonPaste{},
			assert: func(t *testing.T, p *Paste) {
				if p == nil {
					t.Error("paste should not be nil, but was")
				}
			},
		},
		{
			desc: "json paste with data",
			jsonPaste: &jsonPaste{
				ScrapeURL: "https://pastebim.com/github",
				FullURL:   "https://pastebin.com/github",
				Date:      strconv.FormatInt(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Unix(), 10),
				Key:       "testkey",
				Size:      "10",
				Expire:    strconv.FormatInt(time.Date(1970, 2, 1, 0, 0, 0, 0, time.UTC).Unix(), 10),
				Title:     "Golang",
				Syntax:    "go",
				User:      "TwiN",
				Hits:      "3",
			},
			assert: func(t *testing.T, p *Paste) {
				key := "github"
				if p.Key != key {
					t.Errorf("expected key %v; got %v", key, p.Key)
				}
				title := "Golang"
				if p.Title != title {
					t.Errorf("expected title %v; got %v", title, p.Title)
				}
				fullURL := "https://pastebin.com/github"
				if p.URL != fullURL {
					t.Errorf("expected URL %v; got %v", fullURL, p.URL)
				}
				hits := 3
				if p.Hits != hits {
					t.Errorf("expected hits %v; got %v", hits, p.Hits)
				}
				size := 10
				if p.Size != size {
					t.Errorf("expected size %v; got %v", size, p.Size)
				}
				date := "1970-01-01 00:00:00"
				if p.Date.UTC().Format("2006-01-02 15:04:05") != date {
					t.Errorf("expected date %v; got %v", date, p.Date)
				}
				expireDate := "1970-02-01 00:00:00"
				if p.ExpireDate.UTC().Format("2006-01-02 15:04:05") != expireDate {
					t.Errorf("expected expireDate %v; got %v", expireDate, p.ExpireDate)
				}
				visibility := VisibilityPublic
				if p.Visibility != visibility {
					t.Errorf("expected visibility %v; got %v", visibility, p.Visibility)
				}
				syntax := "go"
				if p.Syntax != syntax {
					t.Errorf("expected syntax %v; got %v", syntax, p.Syntax)
				}
				user := "TwiN"
				if p.User != user {
					t.Errorf("expected user %v; got %v", user, p.User)
				}
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := tC.jsonPaste.ToPaste()
			tC.assert(t, p)
		})
	}
}

func TestVisibility(t *testing.T) {
	testCases := []struct {
		desc       string
		visibility Visibility
		want       string
	}{
		{
			desc:       "public visibility",
			visibility: VisibilityPublic,
			want:       "public",
		},
		{
			desc:       "private visibility",
			visibility: VisibilityPrivate,
			want:       "private",
		},
		{
			desc:       "unlisted visibility",
			visibility: VisibilityUnlisted,
			want:       "unlisted",
		},
		{
			desc:       "unknown visibility",
			visibility: Visibility(5355),
			want:       "unknown",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.visibility.String()
			if got != tC.want {
				t.Errorf("expected visibility %s; got %s", tC.want, got)
			}
		})
	}
}
