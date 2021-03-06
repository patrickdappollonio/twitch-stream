package main

import "fmt"

var (
	vlc_path_64 = `C:\Program Files (x86)\VideoLAN\VLC\VLC.exe`
	vlc_path_32 = `C:\Program Files\VideoLAN\VLC\VLC.exe`
)

func getVLCPath() (string, error) {
	var path string

	if exists(vlc_path_32) {
		path = vlc_path_32
	}

	if exists(vlc_path_64) {
		path = vlc_path_64
	}

	if path == "" {
		return "", fmt.Errorf("vlc executable not found!")
	}

	return path, nil
}
