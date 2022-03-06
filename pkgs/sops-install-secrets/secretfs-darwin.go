//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
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

func SecureSymlinkChown(symlinkToCheck, expectedTarget string, owner, group int) error {
	buf := make([]byte, len(expectedTarget)+1) // oversize by one to detect trunc
	n, err := unix.Readlink(symlinkToCheck, buf)
	if err != nil {
		return fmt.Errorf("couldn't readlinkat %s", symlinkToCheck)
	}
	if n > len(expectedTarget) || string(buf[:n]) != expectedTarget {
		return fmt.Errorf("symlink %s does not point to %s", symlinkToCheck, expectedTarget)
	}
	err = unix.Lchown(symlinkToCheck, owner, group)
	if err != nil {
		return fmt.Errorf("cannot change owner of '%s' to %d/%d: %w", symlinkToCheck, owner, group, err)
	}
	return nil
}
