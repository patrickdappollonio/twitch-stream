// Package twitch is a simple library to interact with Twitch
// m3u8 stream URL retrieval in a simple way.
package twitch

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/franela/goreq"
)

const (
	BEST_QUALITY         = "source"
	ACCESS_TOKEN_URL     = "http://api.twitch.tv/api/channels/%s/access_token"
	STREAM_GENERATOR_URL = "http://usher.twitch.tv/api/channel/hls/%s.m3u8?player=twitchweb&token=%s&sig=%s&allow_audio_only=true&allow_source=true&type=any&allow_spectre=false&p=%d"
)

var (
	reQualities  = regexp.MustCompile(`NAME=\"([a-zA-Z_ ]+)\"`)
	reStreamURLs = regexp.MustCompile(`http(|s):\/\/(.*)`)

	ErrCantConnectToTwitch      = errors.New("can't connect to Twitch servers!")
	ErrCantParseAuthResponse    = errors.New("can't parse auth response!")
	ErrNoStreamsFoundForUser    = errors.New("no streams found for the given user")
	ErrCantParseStreamsResponse = errors.New("can't parse the list of streams for the given user")
	ErrStreamQuantitiesMismatch = errors.New("stream qualities / urls mismatch")
	ErrMalformedURLFromAPI      = errors.New("impossible to parse malformed URL from Twitch API")
)

type StreamDetails struct {
	Quality string
	URL     *url.URL
}

// GetTokenAndSignature allows to retrieve a secure token and a given signature
// for a given stream declared as `channel`.
func GetTokenAndSignature(channel string) (string, string, error) {
	// Check if channel is coming
	if strings.TrimSpace(channel) == "" {
		panic("no channel name sent")
	}

	// Perform the request with the given stream
	res, err := goreq.Request{
		Uri: fmt.Sprintf(ACCESS_TOKEN_URL, channel),
	}.Do()

	// Check if there was an error connecting
	// to the twitch website
	if err != nil {
		return "", "", ErrCantConnectToTwitch
	}

	// Close response body once we're done here
	defer res.Body.Close()

	// Create the response's struct based on the json response
	var response struct {
		Token     string `json:"token"`
		Signature string `json:"sig"`
	}

	// Check if there was an error here
	if err := res.Body.FromJsonTo(&response); err != nil {
		return "", "", ErrCantParseAuthResponse
	}

	return response.Token, response.Signature, nil
}

// GetStreamWithQualityOrBest retrieves the whole list of streams and picks
// one based on the quality selected or the best one if the given quality
// was not found.
func GetStreamWithQualityOrBest(channel, token, signature, quality string) (StreamDetails, error) {
	// Check that channel is not empty
	if strings.TrimSpace(channel) == "" {
		panic("no channel name sent")
	}

	// Check that token is not empty
	if strings.TrimSpace(token) == "" {
		panic("no token sent")
	}

	// Check that signature is not empty
	if strings.TrimSpace(signature) == "" {
		panic("no signature sent")
	}

	// Check that quality is not empty
	if strings.TrimSpace(quality) == "" {
		panic("no quality sent")
	}

	// Get all streams
	streams, err := GetAllStreams(channel, token, signature)

	// Check if there was an error
	if err != nil {
		return StreamDetails{}, err
	}

	// Check if there's zero streams retrieved
	if len(streams) == 0 {
		return StreamDetails{}, ErrNoStreamsFoundForUser
	}

	// Iterate over each stream and find both the best and the given one
	// we can safely iterate since at least we have one here
	var stream StreamDetails
	for pos, single := range streams {
		// If this is the quality we're looking for
		// we should stop searching and just continue
		if single.Quality == quality {
			stream = streams[pos]
			break
		}

		// On every sweep we check if this is the best
		// quality available (usually that's for non-partners)
		if single.Quality == BEST_QUALITY {
			stream = streams[pos]
			continue
		}

	}

	return stream, nil
}

// GetAllStreams retrieves all available stream URLs with their given quality.
// For Twitch Partners, they'll have 5 different qualities: source, high, medium,
// low, and mobile. While for the rest, the only available quality is source.
func GetAllStreams(channel, token, signature string) ([]StreamDetails, error) {
	// Check that channel is not empty
	if strings.TrimSpace(channel) == "" {
		panic("no channel name sent")
	}

	// Check that token is not empty
	if strings.TrimSpace(token) == "" {
		panic("no token sent")
	}

	// Check that signature is not empty
	if strings.TrimSpace(signature) == "" {
		panic("no signature sent")
	}

	// Perform the request with the given stream
	res, err := goreq.Request{
		Uri: fmt.Sprintf(STREAM_GENERATOR_URL, channel, url.QueryEscape(token), signature, time.Now().UnixNano()),
	}.Do()

	// Check if there was an error connecting
	// to the twitch website
	if err != nil {
		return []StreamDetails{}, ErrCantConnectToTwitch
	}

	// Close response body once we're done here
	defer res.Body.Close()

	// Convert response to string
	body, err := res.Body.ToString()

	// Check if there was an error trying to
	// convert the body to string
	if err != nil {
		return []StreamDetails{}, ErrCantParseStreamsResponse
	}

	// Find all qualities and urls in the page body
	var (
		availableQualities = reQualities.FindAllString(body, -1)
		availableStreams   = reStreamURLs.FindAllString(body, -1)
	)

	// Check if there are zero values, that means user is not streaming
	if len(availableQualities) == 0 {
		return []StreamDetails{}, ErrNoStreamsFoundForUser
	}

	// Check if we find the same amount of qualities and streams
	if len(availableQualities) != len(availableStreams) {
		return []StreamDetails{}, ErrStreamQuantitiesMismatch
	}

	// Create placeholders
	var errorparsing error
	var streamlist []StreamDetails

	// Iterate over the streams and append them to the list
	for pos := range availableStreams {
		// Parse the given link as a whole URL
		streamlink, err := url.Parse(availableStreams[pos])

		// If it wasn't possible to parse, then stop appending values
		// and break, returning the error
		if err != nil {
			errorparsing = ErrMalformedURLFromAPI
			break
		}

		// Clean quality
		qual := strings.ToLower(availableQualities[pos])
		qual = strings.Replace(qual, `name="`, "", 1)
		qual = strings.Replace(qual, `"`, "", 1)

		// Append the given stream to the list
		streamlist = append(streamlist, StreamDetails{qual, streamlink})
	}

	// Check if there was an error
	if errorparsing != nil {
		return []StreamDetails{}, errorparsing
	}

	// Return the remaining elements
	return streamlist, nil
}
