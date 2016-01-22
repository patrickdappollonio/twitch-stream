package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/patrickdappollonio/twitch-stream/twitch"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// Declare the global app
	app = kingpin.New("twitch-stream", APP_DESCRIPTION)

	// Flags and arguments
	argStreamer = app.Arg(options["streamer"].Name, options["streamer"].Value).Required().String()
	argQuality  = app.Arg(options["quality"].Name, options["quality"].Value).Default("best").String()
	argShowURL  = app.Flag(options["showurl"].Name, options["showurl"].Value).Short(options["showurl"].Short).Bool()
	argDebug    = app.Flag(options["debug"].Name, options["debug"].Value).Short(options["debug"].Short).Bool()
)

func main() {
	// Set app version
	app.Version("2.0.0")

	// Parse all values
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// Placeholders
	var channel, quality string
	var showURL bool

	// Enable or disable debug information
	var buf bytes.Buffer
	logger := log.New(&buf, "[debug] -- ", 0)
	logger.SetOutput(os.Stdout)

	// If no debug, then discard
	if !(*argDebug) {
		logger.SetFlags(0)
		logger.SetOutput(ioutil.Discard)
	}

	// Create a placeholder (to prevent variable shadowing)
	var (
		wasFound     bool
		whichQuality Key
	)

	// Retrieve the quality from the available ones
	whichQuality, wasFound = qualities[*argQuality]

	// Check if the quality exists in the available list
	if !wasFound {
		errorAndExit(fmt.Sprintf("no such quality available: %s", *argQuality))
		return
	}

	// Set variables to placeholders
	channel = fmt.Sprintf("%s", *argStreamer)
	quality = whichQuality.Value
	showURL = *argShowURL

	// Print information to console
	fmt.Printf(LOG_TAG, fmt.Sprintf("acquiring credentials to watch twitch.tv/%s", channel))

	// Try retrieving access details to request later the stream URL
	token, signature, err := twitch.GetTokenAndSignature(channel)

	// If there was an error, print and return
	if err != nil {
		logger.Println(err.Error())
		errorAndExit(fmt.Sprintf("unable to connect to the stream at twitch.tv/%s, try again?", channel))
		return
	}

	// Check if no token is coming, if so, it's likely that the username was wrong
	if strings.TrimSpace(token) == "" {
		errorAndExit(fmt.Sprintf(`no credentials received, check the streamer name at http://www.twitch.tv/%s and try again.`, channel))
		return
	}

	// Printing information message
	fmt.Printf(LOG_TAG, fmt.Sprintf("credentials acquired, retrieving stream URL at %s quality", quality))

	// Retrieve the m3u8 url from Twitch undocumented API
	// with a given quality or best
	stream, err := twitch.GetStreamWithQualityOrBest(channel, token, signature, quality)

	// If there was an error, print and return
	if err != nil {
		logger.Println(err.Error())

		if err == twitch.ErrNoStreamsFoundForUser {
			errorAndExit(fmt.Sprintf(`couldn't find any stream for twitch.tv/%s, is the channel live?`, channel))
			return
		}

		errorAndExit(fmt.Sprintf(`unable to retrieve stream for twitch.tv/%s at "%s" quality`, channel, quality))
		return
	}

	// Check if the stream we choose earlier was the one returned
	if stream.Quality != quality {
		fmt.Printf(LOG_TAG, fmt.Sprintf(`quality "%s" wasn't available, choosing "best"`, quality))
	}

	// If we had to show the URL, we show it and close the app
	if showURL {
		fmt.Printf(LOG_TAG, fmt.Sprintf(`url found for "%s" stream`, channel))
		fmt.Println("URL:", stream.URL)
		fmt.Println(THANKS_MSG)
		os.Exit(0)
		return
	}

	// Try finding the vlc executable
	path, err := getVLCPath()

	// If the app wasn't found, print a message
	if err != nil {
		errorAndExit("vlc app wasn't found on your path!")
		return
	}

	// Inform we're opening vlc
	fmt.Printf(LOG_TAG, fmt.Sprintf(`opening vlc app, please wait until "%s" stream starts playing...`, channel))

	// Try executing the app
	cmd := exec.Command(path, stream.URL.String())
	if err := cmd.Start(); err != nil {
		errorAndExit("impossible to execute vlc app!")
		return
	}

	// Showing exit message
	fmt.Println(THANKS_MSG)
}

func errorAndExit(message string) {
	// Error messages exit with status code != 0
	printWithLog(1, ERROR_TAG, message)
}

func infoAndExit(message string) {
	// Info messages exit with status code 0
	printWithLog(0, LOG_TAG, message)
}

func printWithLog(status int, tag, message string) {
	fmt.Printf(tag, message)
	os.Exit(status)
}
