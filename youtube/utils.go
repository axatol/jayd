package youtube

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	patternValidateVideoID = regexp.MustCompile("^[A-Za-z0-9_-]+$")
)

type URLQuery struct {
	VideoID string
	// TODO
	// PlaylistID string
}

func ParseURL(input string) (*URLQuery, error) {
	normalised := input
	normalised = strings.Replace(normalised, "music.youtube.com", "youtube.com", 1)
	normalised = strings.Replace(normalised, "youtu.be/", "youtube.com/watch?v=", 1)
	normalised = strings.Replace(normalised, "youtube.com/embed/", "youtube.com/watch?v=", 1)
	normalised = strings.Replace(normalised, "/v/", "/watch?v=", 1)
	normalised = strings.Replace(normalised, "/watch#", "/watch?", 1)
	normalised = strings.Replace(normalised, "/playlist", "/watch", 1)
	normalised = strings.Replace(normalised, "youtube.com/shorts/", "youtube.com/watch?v=", 1)

	parsed, err := url.Parse(normalised)
	if err != nil {
		return nil, err
	}

	if !parsed.Query().Has("v") {
		return nil, fmt.Errorf("url did not specify video id: %s", input)
	}

	id := parsed.Query().Get("v")
	if !patternValidateVideoID.MatchString(id) {
		return nil, fmt.Errorf("video id was invalid: %s", id)
	}

	return &URLQuery{VideoID: id}, nil
}
