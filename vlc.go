package main

import (
	"fmt"
	"net/url"
	"os/exec"
)

func openInVLC(stream url.URL) error {
	path, ok := isVLCInstalled()

	if !ok {
		return fmt.Errorf("VideoLAN app not found at %s", path)
	}

	if err := exec.Command(path, stream.String()).Run(); err != nil {
		return err
	}

	return nil
}
