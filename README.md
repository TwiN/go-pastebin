# go-pastebin

![build](https://github.com/TwinProduction/go-pastebin/workflows/build/badge.svg?branch=master) 
[![Go version](https://img.shields.io/github/go-mod/go-version/TwinProduction/go-pastebin.svg)](https://github.com/TwinProduction/go-pastebin)
[![Follow TwinProduction](https://img.shields.io/github/followers/TwinProduction?label=Follow&style=social)](https://github.com/TwinProduction)


[Pastebin.com](https://pastebin.com/) API wrapper for Golang


## Usage
```
go get -u github.com/TwinProduction/go-pastebin
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
| GetRecentPastesUsingScrapingAPI | no          | Retrieves the most recent pastes using Pastebin's scraping API | yes*

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


