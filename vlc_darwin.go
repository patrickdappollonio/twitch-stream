package main

var vlc_path = `/Applications/VLC.app/Contents/MacOS/VLC`

func isVLCInstalled() (string, bool) {
	if fileExists(vlc_path) {
		return vlc_path, true
	}

	return "", false
}
