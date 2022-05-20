package pastebin

import (
	"strconv"
	"strings"
	"time"
)

type xmlPastes struct {
	Pastes []xmlPaste `xml:"paste"`
}

type xmlPaste struct {
	Key         string `xml:"paste_key"`
	Date        int64  `xml:"paste_date"`
	Title       string `xml:"paste_title"`
	Size        int    `xml:"paste_size"`
	ExpireDate  int64  `xml:"paste_expire_date"`
	Private     int    `xml:"paste_private"`
	FormatLong  string `xml:"paste_format_long"`
	FormatShort string `xml:"paste_format_short"`
	URL         string `xml:"paste_url"`
	Hits        int    `xml:"paste_hits"`
}

func (p *xmlPaste) ToPaste(username string) *Paste {
	paste := &Paste{
		Key:        p.Key,
		Title:      p.Title,
		User:       username,
		URL:        p.URL,
		Hits:       p.Hits,
		Size:       p.Size,
		Date:       time.Unix(p.Date, 0),
		ExpireDate: time.Unix(p.ExpireDate, 0),
		Visibility: Visibility(p.Private),
		Syntax:     p.FormatShort,
	}
	return paste
}

type jsonPastes struct {
	Pastes []jsonPaste `json:"pastes"`
}

type jsonPaste struct {
	ScrapeURL string `json:"scrape_url"`
	FullURL   string `json:"full_url"`
	Date      string `json:"date"`
	Key       string `json:"key"`
	Size      string `json:"size"`
	Expire    string `json:"expire"`
	Title     string `json:"title"`
	Syntax    string `json:"syntax"`
	User      string `json:"user"`
	Hits      string `json:"hits"`
}

func (p *jsonPaste) ToPaste() *Paste {
	unixDate, _ := strconv.Atoi(p.Date)
	unixExpire, _ := strconv.Atoi(p.Expire)
	hits, _ := strconv.Atoi(p.Hits)
	size, _ := strconv.Atoi(p.Size)
	paste := &Paste{
		Key:        strings.TrimPrefix(p.FullURL, "https://pastebin.com/"),
		Title:      p.Title,
		URL:        p.FullURL,
		Hits:       hits,
		Size:       size,
		Date:       time.Unix(int64(unixDate), 0),
		ExpireDate: time.Unix(int64(unixExpire), 0),
		Visibility: VisibilityPublic,
		Syntax:     p.Syntax,
		User:       p.User,
	}
	return paste
}

type Paste struct {
	Key        string
	Title      string
	User       string
	URL        string
	Hits       int
	Size       int
	Date       time.Time
	ExpireDate time.Time
	Visibility Visibility
	Syntax     string
}

type Visibility int

const (
	VisibilityPublic   Visibility = 0
	VisibilityUnlisted Visibility = 1
	VisibilityPrivate  Visibility = 2
)

func (v Visibility) String() string {
	switch v {
	case VisibilityPublic:
		return "public"
	case VisibilityUnlisted:
		return "unlisted"
	case VisibilityPrivate:
		return "private"
	default:
		return "unknown"
	}
}

type CreatePasteRequest struct {
	Title      string
	Code       string
	Expiration Expiration

	// Visibility of the paste that will be created.
	// Note that a Client configured without username/password cannot create a private paste
	Visibility Visibility

	// Syntax is the format of the paste (e.g. go, javascript, json, ...)
	// See https://pastebin.com/doc_api#5 for a full list of supported values
	Syntax string
}

// NewCreatePasteRequest creates a new CreatePasteRequest struct
//
// Should be used as parameter to the Client.CreatePaste method
func NewCreatePasteRequest(title, code string, expiration Expiration, visibility Visibility, syntax string) *CreatePasteRequest {
	return &CreatePasteRequest{
		Title:      title,
		Code:       code,
		Expiration: expiration,
		Visibility: visibility,
		Syntax:     syntax,
	}
}

type Expiration string

const (
	ExpirationTenMinutes Expiration = "10M"
	ExpirationOneHour    Expiration = "1H"
	ExpirationOneDay     Expiration = "1D"
	ExpirationOneWeek    Expiration = "1W"
	ExpirationTwoWeeks   Expiration = "2W"
	ExpirationOneMonth   Expiration = "1M"
	ExpirationSixMonth   Expiration = "6M"
	ExpirationOneYear    Expiration = "1Y"
	ExpirationNever      Expiration = "N"
)
