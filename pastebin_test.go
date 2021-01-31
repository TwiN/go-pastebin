package pastebin

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type mockClient struct {
	DoFunc func(request *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(request *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(request)
	}
	return &http.Response{}, nil
}

func init() {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			return nil, nil
		},
	}
}

func TestNewClient(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(request.Body)
			defer request.Body.Close()
			if request.URL.String() == LoginApiUrl && string(body) == "api_dev_key=token&api_user_name=username&api_user_password=password" {
				return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString("session-key"))}, nil
			}
			return &http.Response{
				StatusCode: 403,
				Body:       ioutil.NopCloser(bytes.NewBufferString("forbidden")),
			}, nil
		},
	}
	client, _ := NewClient("username", "password", "token")
	if client.sessionKey != "session-key" {
		t.Errorf("expected %s, got %s", "session-key", client.sessionKey)
	}
}

func TestNewClientWithoutUsernameAndPassword(t *testing.T) {
	client, err := NewClient("", "", "token")
	if err != nil {
		t.Fatal("shouldn't have returned an error, because the only reason an error could be returned is if client.login() was called, but the username was not specified therefore client.login() shouldn't have returned an error")
	}
	if client.developerApiKey != "token" {
		t.Errorf("expected %s, got %s", "token", client.developerApiKey)
	}
}

func TestClient_DeletePaste(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(request.Body)
			defer request.Body.Close()
			if request.URL.String() == LoginApiUrl && string(body) == "api_dev_key=token&api_user_name=username&api_user_password=password" {
				return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString("session-key"))}, nil
			}
			if request.URL.String() == RawApiUrl && string(body) == "api_dev_key=token&api_option=delete&api_paste_key=paste-key&api_user_key=session-key" {
				return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString("deleted"))}, nil
			}
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewBufferString("bad request")),
			}, nil
		},
	}
	client, _ := NewClient("username", "password", "token")
	err := client.DeletePaste("paste-key")
	if err != nil {
		t.Fatal("shouldn't have returned an error, got", err.Error())
	}
}

func TestClient_DeletePasteWhenNotAuthenticated(t *testing.T) {
	client, _ := NewClient("", "", "token")
	err := client.DeletePaste("paste-key")
	if err != ErrNotAuthenticated {
		t.Error("DeletePaste should've returned ErrNotAuthenticated, because only a client configured with a username and password can delete a paste")
	}
}

func TestClient_CreatePaste(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			if request.URL.String() == LoginApiUrl {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString("session-key")),
				}, nil
			}
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString("https://pastebin.com/abcdefgh")),
			}, nil
		},
	}
	client, _ := NewClient("username", "password", "token")
	pasteKey, err := client.CreatePaste(NewCreatePasteRequest("", "", ExpirationTenMinutes, VisibilityPublic, ""))
	if err != nil {
		t.Error("shouldn't have returned an error")
	}
	if pasteKey != "abcdefgh" {
		t.Errorf("expected %s, got %s", "abcdefgh", pasteKey)
	}
}

func TestClient_CreatePasteWithPrivateVisibility(t *testing.T) {
	client, _ := NewClient("", "", "token")
	_, err := client.CreatePaste(NewCreatePasteRequest("", "", ExpirationTenMinutes, VisibilityPrivate, ""))
	if err != ErrNotAuthenticated {
		t.Error("CreatePaste should've returned ErrNotAuthenticated, because only a client configured with a username and password can create a private paste")
	}
}

func TestClient_CreatePasteWhenHTTPRequestReturnsError(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			if request.URL.String() == LoginApiUrl {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString("session-key")),
				}, nil
			}
			return &http.Response{
				StatusCode: 403,
				Body:       ioutil.NopCloser(bytes.NewBufferString("error")),
			}, nil
		},
	}
	client, _ := NewClient("username", "password", "token")
	_, err := client.CreatePaste(NewCreatePasteRequest("", "", ExpirationTenMinutes, VisibilityPublic, ""))
	if err == nil {
		t.Error("should've returned an error")
	}
}

func TestClient_CreatePasteWhenSessionKeyExpired(t *testing.T) {
	numberOfCallsToLoginApiUrl := 0
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			if request.URL.String() == LoginApiUrl {
				numberOfCallsToLoginApiUrl++
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString("session-key")),
				}, nil
			}
			// Mock the behavior of an expired sessionKey, which should trigger an automatic re-login
			if numberOfCallsToLoginApiUrl == 1 {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString("Bad API request, invalid api_user_key")),
				}, nil
			}
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString("https://pastebin.com/abcdefgh")),
			}, nil
		},
	}
	client, _ := NewClient("username", "password", "token")
	pasteKey, err := client.CreatePaste(NewCreatePasteRequest("", "", ExpirationTenMinutes, VisibilityPublic, ""))
	if err != nil {
		t.Fatal("shouldn't have returned an error")
	}
	if pasteKey != "abcdefgh" {
		t.Errorf("expected %s, got %s", "abcdefgh", pasteKey)
	}
	if numberOfCallsToLoginApiUrl != 2 {
		t.Errorf("expected %d calls to LoginApiUrl, got %s", 2, pasteKey)
	}
}

func TestClient_GetUserPasteContent(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			if request.URL.String() == LoginApiUrl {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString("session-key")),
				}, nil
			}
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString("content")),
			}, nil
		},
	}
	client, _ := NewClient("username", "password", "token")
	content, err := client.GetUserPasteContent("does-not-matter-because-client-is-mocked")
	if err != nil {
		t.Error("shouldn't have returned an error")
	}
	if content != "content" {
		t.Errorf("expected %s, got %s", "content", content)
	}
}

func TestClient_GetAllUserPastes(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`<paste>
	<paste_key>fakefake</paste_key>
	<paste_date>1338651885</paste_date>
	<paste_title>Fake Paste</paste_title>
	<paste_size>5555</paste_size>
	<paste_expire_date>0</paste_expire_date>
	<paste_private>1</paste_private>
	<paste_format_long>Go</paste_format_long>
	<paste_format_short>go</paste_format_short>
	<paste_url>https://pastebin.com/fakefake</paste_url>
	<paste_hits>9999</paste_hits>
</paste>`)),
			}, nil
		},
	}
	client, _ := NewClient("username", "password", "token")
	pastes, err := client.GetAllUserPastes()
	if err != nil {
		t.Error("shouldn't have returned an error")
	}
	if len(pastes) != 1 {
		t.Error("should've returned 1 paste, but returned", len(pastes))
	}
	if ExpectedUser := "username"; pastes[0].User != ExpectedUser {
		t.Errorf("expected User to be '%s', got '%s'", ExpectedUser, pastes[0].User)
	}
	if ExpectedKey := "fakefake"; pastes[0].Key != ExpectedKey {
		t.Errorf("expected Key to be '%s', got '%s'", ExpectedKey, pastes[0].Key)
	}
	if ExpectedTitle := "Fake Paste"; pastes[0].Title != ExpectedTitle {
		t.Errorf("expected Title to be '%s', got '%s'", ExpectedTitle, pastes[0].Title)
	}
	if ExpectedSyntax := "go"; pastes[0].Syntax != ExpectedSyntax {
		t.Errorf("expected Syntax to be '%s', got '%s'", ExpectedSyntax, pastes[0].Syntax)
	}
	if ExpectedSize := 5555; pastes[0].Size != ExpectedSize {
		t.Errorf("expected Size to be '%d', got '%d'", ExpectedSize, pastes[0].Size)
	}
	if ExpectedHits := 9999; pastes[0].Hits != ExpectedHits {
		t.Errorf("expected Hits to be '%d', got '%d'", ExpectedHits, pastes[0].Hits)
	}
	if ExpectedVisibility := VisibilityUnlisted; pastes[0].Visibility != ExpectedVisibility {
		t.Errorf("expected Visibility to be '%d', got '%d'", ExpectedVisibility, pastes[0].Visibility)
	}
	if ExpectedDate := int64(1338651885); pastes[0].Date.Unix() != ExpectedDate {
		t.Errorf("expected Date to be '%d', got '%d'", ExpectedDate, pastes[0].Date.Unix())
	}
}

func TestClient_GetAllUserPastesWithoutCredentials(t *testing.T) {
	client, _ := NewClient("", "", "token")
	_, err := client.GetAllUserPastes()
	if err != ErrNotAuthenticated {
		t.Error("should've returned ErrNotAuthenticated, but returned", err)
	}
}

func TestGetPasteContent(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString("this is code")),
			}, nil
		},
	}
	pasteContent, err := GetPasteContent("abcdefgh")
	if err != nil {
		t.Fatal("shouldn't have returned an error, but returned", err)
	}
	if pasteContent != "this is code" {
		t.Errorf("expected '%s', got '%s'", "this is code", pasteContent)
	}
}

func TestGetPasteUsingScrapingAPIWhenPasteKeyInvalid(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString("Error, we cannot find this paste.")),
			}, nil
		},
	}
	_, err := GetPasteUsingScrapingAPI("")
	if ExpectedError := "Error, we cannot find this paste."; err == nil || err.Error() != ExpectedError {
		t.Errorf("error should've been '%s', but was '%s'", ExpectedError, err)
	}
}

func TestGetPasteContentUsingScrapingAPIWhenPasteKeyInvalid(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString("Error, paste key is not valid.")),
			}, nil
		},
	}
	_, err := GetPasteContentUsingScrapingAPI("")
	if ExpectedError := "Error, paste key is not valid."; err == nil || err.Error() != ExpectedError {
		t.Errorf("error should've been '%s', but was '%s'", ExpectedError, err)
	}
}

func TestGetPasteContentUsingScrapingAPIWhenIpBlocked(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 403,
				Body:       ioutil.NopCloser(bytes.NewBufferString("Forbidden: YOUR IP: 1.256.256.256 DOES NOT HAVE ACCESS. VISIT: https://pastebin.com/doc_scraping_api TO GET ACCESS!")),
			}, nil
		},
	}
	_, err := GetPasteContentUsingScrapingAPI("abcdefgh")
	if err == nil {
		t.Error("should've returned an error")
	}
}

func TestGetRecentPastesUsingScrapingAPI(t *testing.T) {
	client = &mockClient{
		DoFunc: func(request *http.Request) (*http.Response, error) {
			if request.URL.String() == ScrapingApiUrl+"?lang=go&limit=1" {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`[{"title":"title"}]`)),
				}, nil
			}
			t.Error(request)
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewBufferString("bad request")),
			}, nil
		},
	}
	pastes, err := GetRecentPastesUsingScrapingAPI("go", 1)
	if err != nil {
		t.Fatal("shouldn't have returned an error, got", err.Error())
	}
	if len(pastes) != 1 {
		t.Error("expected 1 paste to be returned")
	}
}
