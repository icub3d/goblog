// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
)

const (
	version = "0.5.0"
)

func main() {
	// Parse the flags.
	flag.Parse()

	if Version {
		fmt.Println("goblog", version)
		return
	}

	SetupDirectories()

	// First load the templates.
	tmplts, err := LoadTemplates(TemplateDir)
	if err != nil {
		fmt.Println("loading templates:", err)
		os.Exit(1)
	}

	// Next, let's clear out the OutputDir if requested.
	if EmptyOutputDir {
		err = os.RemoveAll(OutputDir)
		if err != nil {
			fmt.Println("cleaning output dir:", err)
			os.Exit(1)
		}
	}

	// Make the output dir.
	err = MakeDirIfNotExists(OutputDir)
	if err != nil {
		fmt.Println("making output dir:", err)
		os.Exit(1)
	}

	// Now, move the static files over.
	err = CopyFilesRecursively(OutputDir, StaticDir)
	if err != nil {
		fmt.Println("making output dir:", err)
		os.Exit(1)
	}

	// Get a list of files from the BlogDir.
	entries, err := GetBlogFiles(BlogDir)
	if err != nil {
		fmt.Println("getting blog file list:", err)
		os.Exit(1)
	}

	// Iteratively Parse each blog for it's useful data and generate a
	// page for each blog.
	for _, blog := range entries {
		contents, err := blog.Parse()
		if err != nil {
			fmt.Println("parsing blog", blog, ":", err)
			os.Exit(1)
		}

		err = tmplts.MakeBlogEntry(OutputDir, blog, contents)
		if err != nil {
			fmt.Println("generating blog html", blog, ":", err)
			os.Exit(1)
		}
	}

	// Generate the about page.
	err = tmplts.MakeAbout(OutputDir)
	if err != nil {
		fmt.Println("generating about.html:", err)
		os.Exit(1)
	}

	// Generate the tags page.
	tags := GetTags(entries)
	err = tmplts.MakeTags(OutputDir, tags.Slice())
	if err != nil {
		fmt.Println("generating tags.html:", err)
		os.Exit(1)
	}

	// Get a sort list of archives.
	dateentries := ParseBlogs(entries)
	a := dateentries.Slice()

	// Generate the archive page.
	err = tmplts.MakeArchive(OutputDir, a)
	if err != nil {
		fmt.Println("generating archive.html:", err)
		os.Exit(1)
	}

	// Generate the index page.
	mostRecent := GetMostRecent(a, MaxIndexEntries)
	err = tmplts.MakeIndex(OutputDir, mostRecent)
	if err != nil {
		fmt.Println("generating index.html:", err)
		os.Exit(1)
	}

	// Generate the RSS feed.
	err = MakeRss(GetMostRecent(a, 10), URL, TemplateDir, OutputDir)
	if err != nil {
		fmt.Println("generating feed.rss:", err)
		fmt.Println("no rss will be available")
	}

}
