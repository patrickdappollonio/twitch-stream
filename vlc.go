package main

import (
	"fmt"
	"net/url"
	"syscall"
)

func openInVLC(stream url.URL) error {
	path, ok := isVLCInstalled()

	if !ok {
		return fmt.Errorf("VideoLAN app not found at %s", path)
	}

	argv := []string{"vlc", stream.String()}

	if _, err := syscall.ForkExec(path, argv, nil); err != nil {
		return err
	}

	return nil
}
