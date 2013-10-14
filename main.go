// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
	"path"
)

const (
	version = "0.5.0"
)

// Version is the flag used to print out the version of the
// application.
var Version bool

// WorkingDir is the directory where that should be prepended to all
// the other configurable directories.
var WorkingDir string

// OutputDir is the directory where the results of the program should
// be written to.
var OutputDir string

// EmptyOutputDir is a flag that determines whether or not the
// OutputDir should be cleaned up before writing file to it.
var EmptyOutputDir bool

// TemplateDir is the diretory where the templates can be found.
var TemplateDir string

// BlogDir is the directory where the blog posts can be found.
var BlogDir string

// StaticDir is the directory where static assests can be found.
var StaticDir string

// URL is the url for this site. The RSS feed will use it to generate links.
var URL string

// MaxIndexEntries is the maximum number of entries to display on the
// index page.
var MaxIndexEntries int

func init() {
	flag.BoolVarP(&Version, "version", "v", false,
		"Output the current version of the application.")

	flag.StringVarP(&WorkingDir, "working-dir", "w", "./",
		"The directory where all the other directories reside. This "+
			"will be prepended to the rest of the configurable directories.")

	flag.StringVarP(&OutputDir, "output-dir", "o", "public",
		"The directory where the results should be placed.")

	flag.BoolVarP(&EmptyOutputDir, "empty-output-dir", "x", false,
		"Before writing to the output-dir, delete anything inside of it.")

	flag.StringVarP(&TemplateDir, "template-dir", "t", "templates",
		"The directory where the site templates are located.")

	flag.StringVarP(&BlogDir, "blog-dir", "b", "blogs",
		"The directory where the blogs are located.")

	flag.StringVarP(&StaticDir, "static-dir", "s", "static",
		"The directory where the static assets are located.")

	flag.StringVarP(&URL, "url", "u", "",
		"The url to be prepended to link in the RSS feed. Defaults to "+
			"the value in the channel <link>.")

	flag.IntVarP(&MaxIndexEntries, "index-entries", "i", 3,
		"The maximum number of entries to display on the index page.")

}

func main() {
	// Parse the flags.
	flag.Parse()

	// Check the version flag first.
	if Version {
		fmt.Println("goblog", version)
		return
	}

	// Setup the directories.
	OutputDir = path.Join(WorkingDir, OutputDir)
	TemplateDir = path.Join(WorkingDir, TemplateDir)
	StaticDir = path.Join(WorkingDir, StaticDir)
	BlogDir = path.Join(WorkingDir, BlogDir)

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

		err = tmplts.MakeEntry(OutputDir, blog, contents)
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
	ebd := GetEntriesByDate(entries)
	err = tmplts.MakeArchive(OutputDir, GetArchives(ebd))
	if err != nil {
		fmt.Println("generating archive.html:", err)
		os.Exit(1)
	}

	// Generate the index page.
	err = tmplts.MakeIndex(OutputDir, ebd[:MaxIndexEntries])
	if err != nil {
		fmt.Println("generating index.html:", err)
		os.Exit(1)
	}

	// Generate the RSS feed.
	err = MakeRss(ebd[:10], URL, TemplateDir, OutputDir)
	if err != nil {
		fmt.Println("generating feed.rss:", err)
		fmt.Println("no rss will be available")
	}
}
