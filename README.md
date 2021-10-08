# go-pastebin

![build](https://github.com/TwiN/go-pastebin/workflows/build/badge.svg?branch=master) 
[![Go Report Card](https://goreportcard.com/badge/github.com/TwiN/go-pastebin)](https://goreportcard.com/report/github.com/TwiN/go-pastebin)
[![codecov](https://codecov.io/gh/TwiN/go-pastebin/branch/master/graph/badge.svg)](https://codecov.io/gh/TwiN/go-pastebin)
[![Go version](https://img.shields.io/github/go-mod/go-version/TwiN/go-pastebin.svg)](https://github.com/TwiN/go-pastebin)
[![Go Reference](https://pkg.go.dev/badge/github.com/TwiN/go-pastebin.svg)](https://pkg.go.dev/github.com/TwiN/go-pastebin)

A [Pastebin.com](https://pastebin.com/) API wrapper in Go.


## Table of Contents

- [Usage](#usage)
  - [Creating a paste](#creating-a-paste)
  - [Deleting a paste](#deleting-a-paste)
  - [Retrieving the content of a paste](#retrieving-the-content-of-a-paste)
    - [GetUserPasteContent](#getuserpastecontent)
    - [GetPasteContent](#getpastecontent)
    - [GetPasteContentUsingScrapingAPI](#getpastecontentusingscrapingapi)
  - [Retrieving paste metadata](#retrieving-paste-metadata)
    - [GetAllUserPastes](#getalluserpastes)
    - [GetPasteUsingScrapingAPI](#getpasteusingscrapingapi)
    - [GetRecentPastesUsingScrapingAPI](#getrecentpastesusingscrapingapi)


## Usage
```
go get -u github.com/TwiN/go-pastebin
```

| Function                        | Client      | Description | PRO          |
|:------------------------------- |:----------- |:----------- |:------------ |
| NewClient                       | n/a         | Creates a new Client | no
| CreatePaste                     | yes         | Creates a new paste and returns the paste key | no
| DeletePaste                     | yes         | Removes a paste that belongs to the authenticated user | no
| GetAllUserPastes                | yes         | Retrieves a list of pastes owned by the authenticated user | no
| GetUserPasteContent             | yes         | Retrieves the content of a paste owned by the authenticated user | no
| GetPasteContent                 | no          | Retrieves the content of a paste using the raw endpoint. This does not require authentication, but only works with public and unlisted pastes. Using this excessively could lead to your IP being blocked. You may want to use GetPasteContentUsingScrapingAPI instead. | no
| GetPasteContentUsingScrapingAPI | no          | Retrieves the content of a paste using Pastebin's scraping API | yes*
| GetPasteUsingScrapingAPI        | no          | Retrieves the metadata of a paste using Pastebin's scraping API | yes*
| GetRecentPastesUsingScrapingAPI | no          | Retrieves a list of recent pastes using Pastebin's scraping API | yes*

\*To use Pastebin's Scraping API, you must [link your IP to your account](https://pastebin.com/doc_scraping_api)

### Creating a paste
You can create a paste by using `pastebin.Client`'s **CreatePaste** function:
```go
client, err := pastebin.NewClient("username", "password", "token")
if err != nil {
	panic(err)
}
pasteKey, err := client.CreatePaste(pastebin.NewCreatePasteRequest("title", "content", pastebin.ExpirationTenMinutes, pastebin.VisibilityUnlisted, "go"))
if err != nil {
	panic(err)
}
fmt.Println("Created paste:", pasteKey)
```
To view the paste on your browser, you can simply append the returned **pasteKey** to `https://pastebin.com/`.

Passing an empty string as username and as password for the client will result in the creation of a guest paste
rather than a paste owned by a user. Note that only authenticated users may create private pastes.


### Deleting a paste
You can delete a paste owned by the user configured in the client by using the **DeletePaste** function:
```go
client, err := pastebin.DeleteClient("username", "password", "token")
if err != nil {
	panic(err)
}
pasteKey, err := client.CreatePaste(pastebin.NewCreatePasteRequest("title", "content", pastebin.ExpirationTenMinutes, pastebin.VisibilityUnlisted, "go"))
if err != nil {
	panic(err)
}
fmt.Println("Created paste:", pasteKey)
```


### Retrieving the content of a paste
There 3 ways to retrieve the content of a paste:

#### GetUserPasteContent
If you own the paste, you should use this.
```go
client, err := pastebin.NewClient("username", "password", "token")
if err != nil {
	panic(err)
}
pasteContent, err := client.GetUserPasteContent("abcdefgh")
if err != nil {
	panic(err)
}
println(pasteContent)
```

#### GetPasteContent
The paste is public or unlisted, and you don't have a Pastebin PRO account.
```go
pasteContent, err := pastebin.GetPasteContent("abcdefgh")
if err != nil {
	panic(err)
}
println(pasteContent)
```
**WARNING:** Using this excessively could lead to your IP being blocked. You may want to use [GetPasteContentUsingScrapingAPI](#getpastecontentusingscrapingapi) instead.

#### GetPasteContentUsingScrapingAPI
The paste is public or unlisted and you have a Pastebin PRO account with your [IP linked](https://pastebin.com/doc_scraping_api).
```go
pasteContent, err := pastebin.GetPasteContentUsingScrapingAPI("abcdefgh")
if err != nil {
	panic(err)
}
println(pasteContent)
```


### Retrieving paste metadata
Just like [retrieving paste content](#retrieving-the-content-of-a-paste), there are many ways to retrieve paste metadata.

<details>
    <summary>List of fields available in paste metadata</summary>

- key
- title
- user
- url
- hits
- size
- date
- expiration date
- visibility (public, unlisted, private)
- syntax
</details>

#### GetAllUserPastes
This will return a list of pastes owned by the user.
```go
client, err := pastebin.NewClient("username", "password", "token")
if err != nil {
	panic(err)
}
pastes, err := client.GetAllUserPastes()
if err != nil {
	panic(err)
}
for _, paste := range pastes {
	fmt.Printf("key=%s title=%s hits=%d visibility=%d url=%s syntax=%s\n", paste.Key, paste.Title, paste.Hits, paste.Visibility, paste.URL, paste.Syntax)
}
```

#### GetPasteUsingScrapingAPI
This will return the metadata of a single paste.
```go
paste, err := pastebin.GetPasteUsingScrapingAPI("abcdefgh")
if err != nil {
	panic(err)
}
fmt.Printf("key=%s title=%s hits=%d visibility=%d url=%s syntax=%s\n", paste.Key, paste.Title, paste.Hits, paste.Visibility, paste.URL, paste.Syntax)
```
Because this function doesn't require authentication, it only works for public and unlisted pastes.

#### GetRecentPastesUsingScrapingAPI
This will return a list of the most recent pastes on Pastebin.
```go
recentPastes, err := pastebin.GetRecentPastesUsingScrapingAPI("", 30)
if err != nil {
	panic(err)
}
for _, paste := range recentPastes { 
    fmt.Printf("key=%s title=%s hits=%d visibility=%d url=%s syntax=%s\n", paste.Key, paste.Title, paste.Hits, paste.Visibility, paste.URL, paste.Syntax)
}
```
This method takes in **syntax** and **limit** as parameters. Leaving the **syntax** string empty applies no filtering. 
The full list of supported values can be found [here](https://pastebin.com/doc_api#5).
