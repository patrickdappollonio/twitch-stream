package main

import "os/exec"

func getVLCPath() (string, error) {
	return exec.LookPath(VLC_APP)
}
