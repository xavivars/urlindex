# UrlIndex-GO

A Go library to store efficiently (both in terms of space and time point of view) URL list.

It could actually be used to store any string-like content, when the goal is checking the existance. 
But it's been designed specifically with URLs in mind, to be used as a library to write 
[Skipper](https://github.com/zalando/skipper) plugins

It's core exposed data structured is 

```go
type UrlIndex struct {
	lastUpdate time.Time
	fst        *vellum.FST
	remotePath string
	localPath  string
	defaultResponse bool
}
```
