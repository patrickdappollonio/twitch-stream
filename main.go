package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"launchpad.net/gnuflag"
)

const (
	AccessTokenURL = "http://api.twitch.tv/api/channels/%s/access_token"
	M3U8StreamURL  = "http://usher.twitch.tv/api/channel/hls/%s.m3u8?player=twitchweb&token=%s&sig=%s&$allow_audio_only=true&allow_source=true&type=any&p=%d"
)

var (
	ErrCantRetrieveAccessToken  = fmt.Errorf("Can't sign up into Twitch servers")
	ErrCantReadM3U8Contents     = fmt.Errorf("Can't retrieve available streams, are you sure the user is still streaming?")
	ErrNoStreamURLFound         = fmt.Errorf("No stream URL found")
	ErrCantRetrieveStreamURLs   = fmt.Errorf("Can't retrieve stream URLs")
	ErrNoStreamAvailableQuality = fmt.Errorf("No streams available for selected quality")
	ErrCantFindStreams          = fmt.Errorf("We couldn't find any stream for the given user")
	ErrForgotParameters         = fmt.Errorf("Please, send both the Twitch username and the desired quality")

	regexQualityMatch = regexp.MustCompile("VIDEO=\".*\"")
	qualities         = []string{"best", "high", "medium", "low", "mobile"}

	strStarting           = "[INFO] Acquiring %s stream data from Twitch Servers\n"
	strOpeningVLC         = "[INFO] Opening %s stream in VLC\n"
	strNoSelectedQuality  = "[INFO] No streams available at %s quality, choosing the best one instead\n"
	strStreamURL          = "[URL] %s\n"
	strBadQualitySelected = "[ERROR] bad quality selected\n\n%s"
	strMissingArgs        = "missing arguments\n\n%s"
)

var help = `twitch-stream - retrieve and play Twitch stream URLs!

USAGE:
	twitch-stream [username] [quality] [global options]
	twitch-stream --help

EXAMPLE:
	twitch-stream patrickdappollonio best
	twitch-stream patrickdappollonio best --vlc

VERSION:
	1.0.0

AVAILABLE QUALITIES:
	best	Original quality used by Twitch (default)
	high
	medium
	low
	mobile

GLOBAL OPTIONS
	--vlc	Open the acquired stream in VLC
	--help	Shows this information
`

type accesstoken struct {
	Token     string `json:"token"`
	Signature string `json:"sig"`
}

type Stream struct {
	Quality string
	URL     url.URL
}

func main() {
	var (
		shouldShowHelp bool
		shouldOpenVLC  bool
	)
	gnuflag.BoolVar(&shouldShowHelp, "help", false, "Should show the Help")
	gnuflag.BoolVar(&shouldOpenVLC, "vlc", false, "Should open the stream in VLC Media Player")
	gnuflag.Parse(true)

	if shouldShowHelp {
		fmt.Println(help)
		os.Exit(0)
		return
	}

	if err := validateNumberOfArguments(os.Args); err != nil {
		errAndExit(err)
		return
	}

	twchannel := strings.ToLower(os.Args[1])
	quality := strings.ToLower(os.Args[2])

	if err := validateQuality(quality); err != nil {
		errAndExit(err)
		return
	}

	fmt.Fprintf(os.Stdout, strStarting, twchannel)

	token, signature, err := getTokenSignature(twchannel)
	if err != nil {
		errAndExit(ErrCantRetrieveAccessToken)
		return
	}

	streamurl, err := getStreamResponse(twchannel, token, signature)
	if err != nil {
		errAndExit(ErrCantReadM3U8Contents)
		return
	}

	strurl, err := selectStreamQualityOrBest(streamurl, quality)
	if err != nil && err != ErrNoStreamAvailableQuality {
		errAndExit(ErrCantFindStreams)
		return
	}

	if shouldOpenVLC {
		fmt.Fprintf(os.Stdout, strOpeningVLC, twchannel)

		if err := openInVLC(strurl); err != nil {
			errAndExit(err)
			return
		}

		os.Exit(0)
		return
	}

	if err == ErrNoStreamAvailableQuality {
		fmt.Fprintf(os.Stdout, strNoSelectedQuality, quality)
	}

	fmt.Fprintf(os.Stdout, strStreamURL, strurl.String())
	os.Exit(0)
}

func errAndExit(err error) {
	fmt.Fprintln(os.Stderr, "[ERROR]", err.Error())
	os.Exit(1)
	return
}

func validateQuality(quality string) error {
	for _, q := range qualities {
		if q == quality {
			return nil
		}
	}

	return fmt.Errorf(strBadQualitySelected, help)
}

func validateNumberOfArguments(args []string) error {
	if len(os.Args) < 3 {
		return fmt.Errorf(strMissingArgs, help)
	}

	return nil
}

func getTokenSignature(twchannel string) (string, string, error) {
	tokenurl := fmt.Sprintf(AccessTokenURL, twchannel)
	response, err := http.Get(tokenurl)

	if err != nil {
		return "", "", err
	}

	defer response.Body.Close()

	var res accesstoken
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return "", "", err
	}

	return res.Token, res.Signature, nil
}

func getStreamResponse(twchannel, token, signature string) ([]Stream, error) {
	m3u8url := fmt.Sprintf(M3U8StreamURL, twchannel, url.QueryEscape(token), signature, rand.Intn(1000))
	var streams []Stream

	response, err := http.Get(m3u8url)

	if err != nil {
		return streams, err
	}

	defer response.Body.Close()

	if response.StatusCode == 404 {
		return streams, ErrNoStreamURLFound
	}

	if response.StatusCode != 200 {
		return streams, ErrCantRetrieveStreamURLs
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return streams, err
	}

	content := string(body[:])

	lines := strings.Split(content, "\n")
	lines = lines[1:]

	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "#EXT-X-STREAM-INF:") && strings.HasPrefix(lines[i+1], "http") {
			quality := regexQualityMatch.FindString(lines[i])
			quality = strings.TrimPrefix(quality, `VIDEO="`)
			quality = strings.TrimSuffix(quality, `"`)
			streamurl, _ := url.Parse(lines[i+1])

			streams = append(streams, Stream{
				Quality: quality,
				URL:     *streamurl,
			})

			i = i + 2
		}
	}

	return streams, nil
}

func selectStreamQualityOrBest(streams []Stream, quality string) (url.URL, error) {
	var selected *Stream

	if quality == "best" {
		quality = "chunked"
	}

	for _, s := range streams {
		if s.Quality == quality {
			return s.URL, nil
		}

		if s.Quality == "chunked" {
			selected = &s
		}
	}

	if selected.Quality != quality {
		return selected.URL, ErrNoStreamAvailableQuality
	}

	return selected.URL, nil
}
