package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"html"
)

const (
	XkcdNotFound = 1 << 0
	XkcdNetworkErr = 1 << 1
	XkcdOtherErr = 1 << 2
	
	PathEnding = "/info.0.json"
)

type XkcdErr struct {
	ErrType int // i.e. XkcdNotFound or XkcdNetworkErr
}

// if comicnumber is 0, we get today's comic
func GetXkcd(comicnumber int) (num int, title string, alt string, imgURL string, err XkcdErr) {
	path := ""
	if comicnumber > 0 {
		path = fmt.Sprintf("https://xkcd.com/%v/info.0.json", comicnumber)
	} else {
		path = "https://xkcd.com/info.0.json"
	}

	resp, e := http.Get(path)
	if e != nil {
		err = XkcdErr{ ErrType: XkcdNetworkErr }
		return
	} else {
		if resp.StatusCode == http.StatusNotFound {
			err = XkcdErr{ ErrType: XkcdNotFound }
			return
		}
	}
	defer resp.Body.Close()
	
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		err = XkcdErr{ ErrType: XkcdOtherErr }
		return
	}
	
	var comicdata map[string]interface{}

	e = json.Unmarshal(body, &comicdata)
	if e != nil {
		err = XkcdErr{ ErrType: XkcdOtherErr }
		return
	}

	num = int(comicdata["num"].(float64))
	title = html.UnescapeString(comicdata["title"].(string))
	alt = html.UnescapeString(comicdata["alt"].(string))
	imgURL = comicdata["img"].(string)

	return
}
