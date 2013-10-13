// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

// MakeDirIfNotExists creates the given directory if it does not exist.
func MakeDirIfNotExists(dir string) error {
	st, err := os.Stat(dir)
	if err != nil {
		// It may just not exist. Make it or error out trying.
		if os.IsNotExist(err) {
			return os.Mkdir(dir, 0750)
		}
	}

	// It exists, make sure it's a directory.
	if !st.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	return nil
}

// CopyFilesRecursively copies the contents of the directory src
// into dest.
func CopyFilesRecursively(dest, src string) error {

	// Read the list of entries for src.
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		s := path.Join(src, file.Name())
		d := path.Join(dest, file.Name())

		// If the file is a directory, make it and then recursively call
		// this function.
		if file.IsDir() {
			err = MakeDirIfNotExists(d)
			if err != nil {
				return err
			}

			err = CopyFilesRecursively(d, s)
			if err != nil {
				return err
			}
		} else {
			// If the file is a file, then copy the file to dest.
			CopyFile(d, s)
		}

		// Set the create/mod times to be the same as the src.
		err = os.Chtimes(d, file.ModTime(), file.ModTime())
		if err != nil {
			return err
		}
	}

	return nil
}

// Copy file makes an exact copy fo the file at src and saves it to
// dest. The contents of dest are overwritten if it exists.
func CopyFile(dest, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

// GetTimes attempts to get the create and last modify times of the
// given path. It first tries using git for the first and last
// commits. If that fails, it uses the system time.
func GetTimes(p string) (time.Time, time.Time, error) {

	// First try to get the git times.
	first, err1 := getGitTime(p, true)
	second, err2 := getGitTime(p, false)

	if err1 == nil && err2 == nil {
		// We got the times, so let's return those!
		return time.Unix(first, 0), time.Unix(second, 0), nil
	}

	// One of them must have failed, so let's get the system times.
	fi, err := os.Stat(p)
	if err != nil {
		return time.Now(), time.Now(), nil
	}

	return fi.ModTime(), fi.ModTime(), nil
}

// getGitTime calls various git commands to get the UNIX Epoch time
// the file was made. if first is true, the original commit time is
// returned, otherwise, the most recent commit time is returned.
func getGitTime(p string, first bool) (int64, error) {
	dir, file := path.Split(p)

	// Get our current directory and defer jumping back to it.
	cwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}
	defer os.Chdir(cwd)

	// Change directories.
	err = os.Chdir(dir)
	if err != nil {
		return 0, err
	}

	// Call the revlist
	cmd := exec.Command("git", "rev-list", "--max-parents=1", "HEAD", file)
	revs, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Get the first or last item.
	revslist := strings.Split(string(revs), "\n")
	var rev string
	if !first {
		rev = revslist[0]
	} else {
		// The last one will actually be an empty line, so we want the
		// second to last.
		if len(revslist) < 2 {
			rev = revslist[0]
		} else {
			rev = revslist[len(revslist)-2]
		}
	}

	if rev == "" {
		return 0, fmt.Errorf("not found")
	}

	// Get the commit time.
	cmd = exec.Command("git", "show", "-s", "--format=%at", rev)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Convert to int.
	i, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil

}
