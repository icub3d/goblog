// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	flag "github.com/ogier/pflag"
	"path"
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
		"The url to be prepended to link in the RSS feed. Defaults to the value in the channel <link>.")

	flag.IntVarP(&MaxIndexEntries, "index-entries", "i", 3,
		"The maximum number of entries to display on the index page.")

}

// SetupDirectories is a helper function that prepends the working
// directory to the other directories.
func SetupDirectories() {
	OutputDir = path.Join(WorkingDir, OutputDir)
	TemplateDir = path.Join(WorkingDir, TemplateDir)
	StaticDir = path.Join(WorkingDir, StaticDir)
	BlogDir = path.Join(WorkingDir, BlogDir)
}
