//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"os"
)

func MountSecretFs(mountpoint string, keysGid int) error {
	if err := os.MkdirAll(mountpoint, 0751); err != nil {
		return fmt.Errorf("Cannot create directory '%s': %w", mountpoint, err)
	}

	if err := os.Chown(mountpoint, 0, int(keysGid)); err != nil {
		return fmt.Errorf("Cannot change owner/group of '%s' to 0/%d: %w", mountpoint, keysGid, err)
	}

	return nil
}
