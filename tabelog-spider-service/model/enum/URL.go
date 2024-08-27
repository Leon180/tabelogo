package enum

import "net/url"

type URL string

func (url URL) ToString() string {
	return string(url)
}

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

const (
	TabelogBaseURL URL = "https://tabelog.com/"
)

type SubDirectory string

func (s SubDirectory) ToString() string {
	return string(s)
}

const (
	TabelogSubDirectoryRstLst      SubDirectory = "rstLst"
	TabelogSubDirectoryRstPhotoLst SubDirectory = "dtlphotolst"
)

type Param string

func (p Param) ToString() string {
	return string(p)
}

const (
	TabelogParamVS Param = "vs"
	TabelogParamSK Param = "sk"
	TabelogParamSW Param = "sw"
)
