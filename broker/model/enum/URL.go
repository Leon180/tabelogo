package enum

import "net/url"

type URL string

func (url URL) ToString() string {
	return string(url)
}

const (
	GoogleMapPlaceV1URL         URL = "https://places.googleapis.com/v1/places/"
	GoogleMapPlaceSearchTextURL URL = "https://places.googleapis.com/v1/places:searchText"
)

func (url URL) AddSubDirectory(subDirectory string) URL {
	if len(url) == 0 {
		return URL(subDirectory)
	}
	if string(url)[len(url)-1] == '/' {
		return URL(url.ToString() + subDirectory)
	}
	return URL(url.ToString() + "/" + subDirectory)
}

func (b URL) AddParams(params map[Param]string) URL {
	newURL := string(b) + "?"
	for k, v := range params {
		newURL = newURL + k.ToString() + "=" + url.QueryEscape(v) + "&"
	}
	if newURL[len(newURL)-1] == '&' {
		newURL = newURL[:len(newURL)-1]
	}
	return URL(newURL)
}

type SubDirectory string

func (s SubDirectory) ToString() string {
	return string(s)
}

type Param string

func (p Param) ToString() string {
	return string(p)
}

const (
	APIMask      Param = "fields"
	Key          Param = "key"
	LanguageCode Param = "languageCode"
)
