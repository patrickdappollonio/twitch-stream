package main

import "fmt"

var vlc_path = `/Applications/VLC.app/Contents/MacOS/VLC`

func getVLCPath() (string, error) {
	if exists(vlc_path) {
		return vlc_path, nil
	}

	return "", fmt.Errorf("vlc executable not found at %s", vlc_path)
}
