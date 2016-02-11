package main

type Key struct {
	Name, Value string
	Short       byte
}

const (
	APP_DESCRIPTION = `twitch-stream is a cross platform application inspired by
livestreamer but focused only on Twitch. It allows you to select a stream you
like and play it on VLC by default (or get the URL of the stream so you can play
it on any application on your own by passing the parameter "-u") by opening a new
window. If you want to have multiple streams open, just execute the application
several times with different streams.

Found a problem with the app? I would love to know more. Please fill an issue
in the Github page at https://github.com/patrickdappollonio/twitch-stream/issues

This application was created by Patrick D'appollonio, which is
http://www.twitch.tv/patrickdappollonio on Twitch. Follow me!`
	LOG_TAG    = "[info] -- %s \n"
	ERROR_TAG  = "[error] -- %s \n"
	VLC_APP    = "vlc"
	THANKS_MSG = "—————— Thanks for using twitch-stream! ——————"
)

var (
	qualities = map[string]Key{
		"best":   Key{Name: "best", Value: "source"},
		"high":   Key{Name: "high", Value: "high"},
		"medium": Key{Name: "medium", Value: "medium"},
		"low":    Key{Name: "low", Value: "low"},
		"mobile": Key{Name: "mobile", Value: "mobile"},
		"audio":  Key{Name: "audio", Value: "audio only"},
	}

	options = map[string]Key{
		"streamer": Key{Name: "streamer", Value: "The twitch username of the stream you want to watch."},
		"quality":  Key{Name: "quality", Value: `The quality of the stream you want to see. Available qualities are "best", "high", "medium", "low", "mobile" and "audio".`},
		"showurl":  Key{Name: "url", Value: "If declared, it will print the stream URL to the console instead.", Short: 'u'},
		"debug":    Key{Name: "debug", Value: "If declared, debug information will be printed to the console.", Short: 'd'},
	}
)
