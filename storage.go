package main

type Storage interface {
	Shorten(Url string, exp int64) (string, error)
	ShortenLinkInfo(eid string) (interface{}, error)
	Unshorten(eid string) (string, error)
}
