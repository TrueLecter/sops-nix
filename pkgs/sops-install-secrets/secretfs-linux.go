//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

const RAMFS_MAGIC int32 = -2054924042

func MountSecretFs(mountpoint string, keysGid int) error {
	if err := os.MkdirAll(mountpoint, 0751); err != nil {
		return fmt.Errorf("Cannot create directory '%s': %w", mountpoint, err)
	}

	buf := unix.Statfs_t{}
	if err := unix.Statfs(mountpoint, &buf); err != nil {
		return fmt.Errorf("Cannot get statfs for directory '%s': %w", mountpoint, err)
	}
	if int32(buf.Type) != RAMFS_MAGIC {
		if err := unix.Mount("none", mountpoint, "ramfs", unix.MS_NODEV|unix.MS_NOSUID, "mode=0751"); err != nil {
			return fmt.Errorf("Cannot mount: %s", err)
		}
	}

	if err := os.Chown(mountpoint, 0, int(keysGid)); err != nil {
		return fmt.Errorf("Cannot change owner/group of '%s' to 0/%d: %w", mountpoint, keysGid, err)
	}

	return nil
}
