package pastebin

import "testing"

func TestNewClient(t *testing.T) {
	client, err := NewClient("", "", "DEV_KEY_HERE")
	if err != nil {
		t.Fatal("Shouldn't have returned an error, because the only reason an error could be returned is if client.login() was called, but the username was not specified therefore client.login() shouldn't have returned an error")
	}
	if client.developerApiKey != "DEV_KEY_HERE" {
		t.Errorf("expected %s, got %s", "DEV_KEY_HERE", client.developerApiKey)
	}
}

func TestClient_DeletePaste(t *testing.T) {
	client, _ := NewClient("", "", "DEV_KEY_HERE")
	err := client.DeletePaste("paste-key")
	if err != ErrNotAuthenticated {
		t.Error("DeletePaste should've instantly returned ErrNotAuthenticated, because only a client configured with a username and password can delete a paste")
	}
}

func TestClient_CreatePasteWithPrivateVisibility(t *testing.T) {
	client, _ := NewClient("", "", "DEV_KEY_HERE")
	_, err := client.CreatePaste(NewCreatePasteRequest("", "", ExpirationTenMinutes, VisibilityPrivate, ""))
	if err != ErrNotAuthenticated {
		t.Error("CreatePaste should've instantly returned ErrNotAuthenticated, because only a client configured with a username and password can create a private paste")
	}
}
